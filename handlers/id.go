package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("id", handleID)
}

// handleID handles the id command.
// All flags are boolean and kept verbatim.
// The optional positional argument (username or UID) collapses to <user>.
func handleID(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			continue
		}

		// Positional: username or UID
		result = append(result, "<user>")
	}

	result = append(result, redirects...)
	return result
}
