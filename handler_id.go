package shellshape

func init() {
	Register("id", handleID)
}

// handleID handles the id command.
// All flags are boolean and kept verbatim.
// The optional positional argument (username or UID) collapses to <user>.
func handleID(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	for _, tok := range args {
		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			continue
		}

		// Positional: username or UID
		result = append(result, "<user>")
	}

	result = append(result, redirects...)
	return result
}
