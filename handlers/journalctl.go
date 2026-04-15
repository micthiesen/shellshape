package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("journalctl", handleJournalctl)
}

// journalctl field match: _FIELD=value (uppercase field name starting with _).
func isJournalFieldMatch(tok string) bool {
	if !strings.HasPrefix(tok, "_") {
		return false
	}
	idx := strings.Index(tok, "=")
	if idx < 2 {
		return false
	}
	field := tok[1:idx]
	for _, c := range field {
		if c < 'A' || c > 'Z' {
			if c != '_' {
				return false
			}
		}
	}
	return true
}

// handleJournalctl handles journalctl commands.
// Value flags (-u, -p, -t, -o, --since, --until, etc.) collapse their argument.
// Numeric flags (-n) collapse to N.
// Path flags (-D, --directory, --file) collapse to <path>.
// Grep flag (-g) collapses to <pattern>.
// Field matches (_PID=123, _COMM=sshd) collapse to <field-match>.
// Executable paths use classifyToken.
func handleJournalctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-u": true, "--unit": true,
		"-t": true, "--identifier": true,
		"-p": true, "--priority": true,
		"-o": true, "--output": true,
		"-b": true, "--boot": true,
		"--since": true, "--until": true,
		"--cursor": true, "--after-cursor": true,
		"-F": true, "--field": true,
		"--vacuum-size": true, "--vacuum-time": true, "--vacuum-files": true,
		"--case-sensitive": true,
		"--facility":       true,
		"--root":           true,
	}

	numericFlags := map[string]bool{
		"-n": true, "--lines": true,
	}

	pathFlags := map[string]bool{
		"-D": true, "--directory": true,
		"--file": true,
	}

	grepFlags := map[string]bool{
		"-g": true, "--grep": true,
	}

	// Flags that can be boolean (no argument) or take an argument.
	// -b can be used alone or with a boot offset/ID.
	optionalArgFlags := map[string]bool{
		"-b": true, "--boot": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: grepFlags, Placeholder: "<pattern>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			// For optional-arg flags, only consume next token if it looks like a value
			if optionalArgFlags[tok] {
				if i < len(args) && !shellshape.IsFlagToken(args[i]) {
					result = shellshape.EmitPositional(result, args[i], "<val>")
					i++
				}
			} else if i < len(args) {
				result = shellshape.EmitPositional(result, args[i], "<val>")
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: field match or path
		if isJournalFieldMatch(tok) {
			result = append(result, "<field-match>")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
