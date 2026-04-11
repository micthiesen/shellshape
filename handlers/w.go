package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("w", handleW)
}

// handleW handles the w command (show who is logged in and what they are doing).
// All flags are boolean and preserved verbatim.
// Positional arguments are user names; they collapse to a single <user>.
func handleW(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	hasUser := false

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			if hasUser {
				result = append(result, "<user>")
				hasUser = false
			}
			result = append(result, tok)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// Positional: user name, collapse all to single <user>.
		hasUser = true
	}
	if hasUser {
		result = append(result, "<user>")
	}

	result = append(result, redirects...)
	return result
}
