package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("shellshape", handleShellshape)
}

// handleShellshape handles the shellshape command itself.
// All positional arguments are the input command being analyzed,
// so they collapse to a single <command> placeholder. Subshells are preserved.
func handleShellshape(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	hasPositional := false

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			if hasPositional {
				result = append(result, "<command>")
				hasPositional = false
			}
			result = append(result, tok)
			continue
		}
		hasPositional = true
	}
	if hasPositional {
		result = append(result, "<command>")
	}

	result = append(result, redirects...)
	return result
}
