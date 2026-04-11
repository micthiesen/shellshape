package shellshape

import (
	"regexp"
	"strings"
)

func init() {
	Register("git", handleGit, HandlerOptions{HasSubcommands: true})
}

var stashRefRE = regexp.MustCompile(`^stash@\{\d+\}$`)

// handleGit handles git subcommand arguments.
// Since git is in subcommandExecutables, the subcommand (commit, push, log, etc.)
// is already consumed before this handler is called.
//
// Message flags (-m, --message) collapse their argument to <str>.
// Value flags (--author, --format, --since, --onto, etc.) collapse to <val>.
// Numeric flags (-n, --depth, -U) collapse to N.
// Path flags (-C, -f, --pathspec-from-file) collapse to <path>.
// Global config flags (-c key=value) collapse to <val>.
// Boolean flags are kept verbatim.
// Stash references (stash@{N}) normalize the index to N.
// Remaining positionals use classifyToken.
func handleGit(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Long flags that can appear with = and should collapse the value.
	eqValFlags := map[string]bool{
		"--pretty": true,
		"--format": true,
	}

	// Flags whose next token is a message string.
	msgFlags := map[string]bool{
		"-m": true, "--message": true,
	}

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"--author": true,
		"--date":   true,
		"--format": true,
		"--pretty": true,
		"--since":  true, "--after": true,
		"--until": true, "--before": true,
		"--grep":            true,
		"--diff-filter":     true,
		"--onto":            true,
		"--set-upstream-to": true,
		"--strategy":        true, "-s": true,
		"--strategy-option": true, "-X": true,
		"--remote":  true,
		"-c":        true,
		"--cleanup": true,
		"-C":        true,
		"--fixup":   true,
		"--squash":  true,
		"--trailer": true,
		"--output":  true, "-o": true,
	}

	// Flags whose next token is a number.
	numericFlags := map[string]bool{
		"-n": true, "--max-count": true,
		"--depth": true, "--deepen": true,
		"-U": true, "--unified": true,
		"-j": true, "--jobs": true,
		"--skip":   true,
		"--abbrev": true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"--pathspec-from-file": true,
		"--work-tree":          true,
		"--git-dir":            true,
		"--template":           true,
	}

	// Fused short flags where -m is combined with other flags (e.g., -am).
	// When we see -am, we know -m consumes the next token as a message.
	isFusedMsgFlag := func(tok string) bool {
		if len(tok) < 3 || tok[0] != '-' || tok[1] == '-' {
			return false
		}
		// -am, -sm, etc. — ends with 'm'
		return tok[len(tok)-1] == 'm'
	}

	var result []string
	configKeysSeen := 0
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Stash references: stash@{0}, stash@{2} → stash@{N}
		if stashRefRE.MatchString(tok) {
			result = append(result, "stash@{N}")
			i++
			continue
		}

		if msgFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<str>")
				}
				i++
			}
			continue
		}

		// Fused message flags like -am
		if isFusedMsgFlag(tok) {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<str>")
				}
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		// Handle --flag=value forms for known flags.
		// Use classifyToken for consistent treatment with the generic path.
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			eqIdx := strings.Index(tok, "=")
			flagName := tok[:eqIdx]
			if eqValFlags[flagName] || valFlags[flagName] || msgFlags[flagName] || pathFlags[flagName] || numericFlags[flagName] {
				result = append(result, flagName+"=<val>")
				i++
				continue
			}
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Config subcommand: second positional (after the key) is a value.
		if subcommand == "config" && configKeysSeen > 0 {
			result = append(result, "<val>")
			i++
			continue
		}

		// Positional: classify generically
		classified := classifyToken(tok)
		if subcommand == "config" && classified == "<dotted-id>" {
			configKeysSeen++
		}
		result = append(result, classified)
		i++
	}

	result = append(result, redirects...)
	return result
}
