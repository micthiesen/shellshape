package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("task", handleTask)
}

// handleTask handles the go-task (taskfile.dev) command.
// Task names are kept verbatim since they define the command shape.
// Variables (KEY=VALUE positionals) collapse to <var>=<val>.
// Everything after "--" is CLI_ARGS passed to tasks, collapsed to <str>.
// Path flags (-d, -t) collapse to <path>, numeric flags (-C) to N,
// and value flags (-o, -I, -c, --sort, etc.) to <val>.
func handleTask(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-d": true, "--dir": true,
		"-t": true, "--taskfile": true,
	}

	numericFlags := map[string]bool{
		"-C": true, "--concurrency": true,
	}

	valFlags := map[string]bool{
		"-o": true, "--output": true,
		"-I": true, "--interval": true,
		"-c": true, "--color": true,
		"--sort":               true,
		"--output-group-begin": true,
		"--output-group-end":   true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		// "--" separator: everything after is CLI_ARGS passed to the task.
		// Collapse all passthrough args to <str> (they're opaque data).
		if tok == "--" {
			result = append(result, "--")
			i++
			for i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<str>")
				}
				i++
			}
			break
		}

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused long flags (--flag=value)
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: check for KEY=VALUE variable syntax
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 && isTaskVariable(tok[:eqIdx]) {
			result = append(result, "<var>=<val>")
			i++
			continue
		}

		// Task name: keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}

// isTaskVariable checks if s looks like a variable name (letters, digits, underscores,
// starting with a letter or underscore).
func isTaskVariable(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, c := range s {
		if i == 0 {
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_') {
				return false
			}
		} else {
			if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
				return false
			}
		}
	}
	return true
}
