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
		return collapseRepeatedPlaceholders(normalizeSegment(parts[0].text))
	}

	var rendered []string
	lastWasSegment := false
	for _, part := range parts {
		seg := normalizeSegment(strings.TrimSpace(part.text))
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
	out = collapseRepeatedPlaceholders(out)
	out = collapseRepeatedSegments(out)
	return out
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
	for i < len(tokens) && IsEnvAssignment(tokens[i]) {
		i++
	}
	if i < len(tokens) {
		return tokens[i]
	}
	return ""
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

// normalizeSegment routes a single segment to either compound or simple normalization.
func normalizeSegment(seg string) string {
	seg = strings.TrimSpace(seg)
	if isCompoundCommand(seg) {
		return normalizeCompound(seg)
	}
	return normalizeSingleCommand(seg)
}

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
	for i < len(tokens) && IsEnvAssignment(tokens[i]) {
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
	if hasSubcommands(exe) && i < len(tokens) {
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
		if RedirectConsumeNext[tok] && i+1 < len(tokens) {
			placeholder := "<path>"
			if tok == "<<<" {
				placeholder = "<str>"
			}
			result = append(result, tok, placeholder)
			i += 2
			continue
		}
		if RedirectStandalone[tok] {
			result = append(result, tok)
			i++
			continue
		}
		if interpWithCode && codeFlags[tok] && i+1 < len(tokens) {
			result = append(result, tok, "<code>")
			i += 2
			continue
		}
		result = append(result, ClassifyToken(tok))
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

// collapseRepeatedSegments collapses consecutive identical segments separated
// by ";" into a single instance with ";+". Only ";" is collapsed since "&&" and
// "||" imply ordering dependency that is semantically meaningful.
func collapseRepeatedSegments(shape string) string {
	// Split into tokens of (segment, operator) pairs.
	// We only collapse on ";", so split on " ; " boundaries.
	const sep = " ; "
	if !strings.Contains(shape, sep) {
		return shape
	}

	// Split on " ; " only — leave && and || intact within segments.
	parts := splitOnSemicolon(shape)
	if len(parts) <= 1 {
		return shape
	}

	// Find consecutive runs of identical segments and collapse them.
	type segment struct {
		text      string
		collapsed bool
	}
	var segments []segment
	i := 0
	for i < len(parts) {
		j := i + 1
		for j < len(parts) && parts[j] == parts[i] {
			j++
		}
		segments = append(segments, segment{parts[i], j-i >= 2})
		i = j
	}

	// Rebuild the shape string.
	var b strings.Builder
	for idx, seg := range segments {
		if idx > 0 {
			// After a collapsed segment, just a space (the ";+" acts as separator).
			// Between non-collapsed segments, restore " ; ".
			if segments[idx-1].collapsed {
				b.WriteString(" ")
			} else {
				b.WriteString(" ; ")
			}
		}
		b.WriteString(seg.text)
		if seg.collapsed {
			b.WriteString(" ;+")
		}
	}
	return b.String()
}

// splitOnSemicolon splits a shape string on " ; " boundaries, but only at
// the top level (not inside segments that contain && or ||).
func splitOnSemicolon(shape string) []string {
	var parts []string
	rest := shape
	for {
		idx := strings.Index(rest, " ; ")
		if idx < 0 {
			parts = append(parts, rest)
			break
		}
		parts = append(parts, rest[:idx])
		rest = rest[idx+len(" ; "):]
	}
	return parts
}

func fallbackNormalize(seg string) string {
	words := strings.Fields(seg)
	if len(words) == 0 {
		return ""
	}
	result := []string{words[0]}
	for _, w := range words[1:] {
		result = append(result, ClassifyToken(w))
	}
	return strings.Join(result, " ")
}

func collapseWhitespace(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}
