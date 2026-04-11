package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("wc", handleWc)
}

// handleWc handles the wc command.
// All flags are boolean (no flag takes an argument).
// All positional arguments are file paths.
func handleWc(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// All positionals are file paths
		result = append(result, "<path>")
	}

	result = append(result, redirects...)
	return result
}
