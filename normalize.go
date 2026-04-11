package shellshape

import (
	"regexp"
	"strings"
)

// Normalize returns a stable shape for the given bash command.
func Normalize(command string) string {
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}

	command = joinContinuations(command)
	command = collapseHeredocs(command)
	command = stripComments(command)
	command = strings.TrimSpace(command)
	if command == "" {
		return ""
	}

	parts := splitTopLevel(command, splitOps)

	if len(parts) == 1 && parts[0].sep == "" {
		return collapseRepeatedPlaceholders(normalizeSingleCommand(parts[0].text))
	}

	var rendered []string
	lastWasSegment := false
	for _, part := range parts {
		seg := normalizeSingleCommand(strings.TrimSpace(part.text))
		if seg != "" {
			rendered = append(rendered, seg)
			lastWasSegment = true
		}
		if part.sep != "" && lastWasSegment {
			rendered = append(rendered, part.sep)
			lastWasSegment = false
		}
	}

	// Drop trailing separator.
	if len(rendered) > 0 {
		last := rendered[len(rendered)-1]
		for _, op := range splitOps {
			if last == op || last == ";" {
				rendered = rendered[:len(rendered)-1]
				break
			}
		}
	}

	out := strings.Join(rendered, " ")
	out = strings.TrimSpace(out)
	out = collapseWhitespace(out)
	return collapseRepeatedPlaceholders(out)
}

// ExecutableOf returns the leading executable name from a shape string.
func ExecutableOf(shape string) string {
	if shape == "" {
		return ""
	}
	// Take the first segment (before any operator).
	for _, op := range splitOps {
		needle := " " + op + " "
		if idx := strings.Index(shape, needle); idx != -1 {
			shape = shape[:idx]
		}
	}
	tokens := strings.Fields(shape)
	i := 0
	for i < len(tokens) && isEnvAssignment(tokens[i]) {
		i++
	}
	if i < len(tokens) {
		return tokens[i]
	}
	return ""
}

// Executables whose second positional token is a subcommand.
var subcommandExecutables = map[string]bool{
	"git": true, "gh": true, "pnpm": true, "npm": true, "yarn": true,
	"bun": true, "cargo": true, "go": true, "docker": true,
	"kubectl": true, "aws": true, "gcloud": true, "az": true,
	"terraform": true, "helm": true, "make": true,
	"pip": true, "pip3": true, "uv": true, "poetry": true, "brew": true,
	"apt": true, "dnf": true, "pacman": true, "snap": true,
	"systemctl": true, "launchctl": true, "sst": true, "pulumi": true,
	"ollama": true, "lms": true, "fnm": true,
	"tmux": true,
}

// Interpreters that accept inline code via -c / -e / --command.
var interpretersWithCodeFlag = map[string]bool{
	"python": true, "python3": true, "node": true, "deno": true,
	"bash": true, "sh": true, "zsh": true,
	"ruby": true, "perl": true, "lua": true, "php": true,
}

var codeFlags = map[string]bool{"-c": true, "-e": true, "--command": true}

// Shells used as script runners.
var shellScriptRunners = map[string]bool{"bash": true, "sh": true, "zsh": true}

func normalizeSingleCommand(seg string) string {
	seg = strings.TrimSpace(seg)
	if seg == "" {
		return ""
	}

	// Recursively normalize $(...) substitutions.
	seg = normalizeSubstitutions(seg)

	tokens, err := shelxSplit(seg)
	if err != nil {
		return fallbackNormalize(seg)
	}
	if len(tokens) == 0 {
		return ""
	}

	// Strip leading env var assignments.
	var envParts []string
	i := 0
	for i < len(tokens) && isEnvAssignment(tokens[i]) {
		name := tokens[i][:strings.Index(tokens[i], "=")+1]
		envParts = append(envParts, name+"<val>")
		i++
	}

	if i >= len(tokens) {
		return strings.Join(envParts, " ")
	}

	exe := tokens[i]
	i++

	// Absolute-path or tilde-path executables: strip to basename.
	// Relative paths (./script.sh, ../bin/run) are left alone.
	if (strings.HasPrefix(exe, "/") || strings.HasPrefix(exe, "~")) && strings.Contains(exe, "/") {
		parts := strings.Split(exe, "/")
		exe = parts[len(parts)-1]
	}

	// Shell-as-script-runner.
	if shellScriptRunners[exe] && i < len(tokens) {
		firstArg := tokens[i]
		if !strings.HasPrefix(firstArg, "-") && strings.Contains(firstArg, "/") {
			parts := strings.Split(firstArg, "/")
			exe = parts[len(parts)-1]
			i++
		} else if !strings.HasPrefix(firstArg, "-") && strings.Contains(firstArg, ".") {
			exe = firstArg
			i++
		}
	}

	result := append(envParts, exe)

	var subcommand string
	// Subcommand detection.
	if subcommandExecutables[exe] && i < len(tokens) {
		subcommand = tokens[i]
		result = append(result, tokens[i])
		i++
	}

	// Per-executable handler.
	if handler, ok := handlers[exe]; ok {
		result = append(result, handler(subcommand, tokens[i:])...)
		return strings.Join(result, " ")
	}

	// Generic per-token classification.
	interpWithCode := interpretersWithCodeFlag[exe]
	for i < len(tokens) {
		tok := tokens[i]
		if redirectConsumeNext[tok] && i+1 < len(tokens) {
			result = append(result, tok, "<path>")
			i += 2
			continue
		}
		if redirectStandalone[tok] {
			result = append(result, tok)
			i++
			continue
		}
		if interpWithCode && codeFlags[tok] && i+1 < len(tokens) {
			result = append(result, tok, "<code>")
			i += 2
			continue
		}
		result = append(result, classifyToken(tok))
		i++
	}

	return strings.Join(result, " ")
}

var subshellRE = regexp.MustCompile(`\$\(`)
var backtickRE = regexp.MustCompile("`[^`]*`")

func normalizeSubstitutions(s string) string {
	// Backticks first: replace with single-quoted opaque placeholder.
	s = backtickRE.ReplaceAllString(s, "'$(<subshell>)'")

	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '$' && i+1 < len(s) && s[i+1] == '(' {
			// Find matching close paren with depth tracking.
			depth := 1
			j := i + 2
			inSingle := false
			inDouble := false
			for j < len(s) && depth > 0 {
				c := s[j]
				if inSingle {
					if c == '\'' {
						inSingle = false
					}
				} else if inDouble {
					if c == '"' {
						inDouble = false
					}
				} else {
					switch c {
					case '\'':
						inSingle = true
					case '"':
						inDouble = true
					case '(':
						depth++
					case ')':
						depth--
						if depth == 0 {
							break
						}
					}
				}
				if depth > 0 {
					j++
				}
			}
			if depth == 0 {
				inner := s[i+2 : j]
				normalizedInner := Normalize(inner)
				// Single-quote-wrap so shlex keeps the subshell atomic.
				out.WriteString("'$(")
				out.WriteString(normalizedInner)
				out.WriteString(")'")
				i = j + 1
				continue
			}
			// Unbalanced: treat rest as opaque.
			out.WriteString(s[i:])
			return out.String()
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

// collapseRepeatedPlaceholders collapses runs of 2+ identical <...> placeholders.
var collapsibleRE = regexp.MustCompile(`^<[a-z][-a-z]*>$`)

func collapseRepeatedPlaceholders(shape string) string {
	if shape == "" {
		return shape
	}
	tokens := strings.Split(shape, " ")
	var out []string
	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		if collapsibleRE.MatchString(tok) {
			j := i + 1
			for j < len(tokens) && tokens[j] == tok {
				j++
			}
			if j-i >= 2 {
				out = append(out, tok+"+")
			} else {
				out = append(out, tok)
			}
			i = j
		} else {
			out = append(out, tok)
			i++
		}
	}
	return strings.Join(out, " ")
}

func fallbackNormalize(seg string) string {
	words := strings.Fields(seg)
	if len(words) == 0 {
		return ""
	}
	result := []string{words[0]}
	for _, w := range words[1:] {
		result = append(result, classifyToken(w))
	}
	return strings.Join(result, " ")
}

func collapseWhitespace(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}
