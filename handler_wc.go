package shellshape

func init() {
	Register("wc", handleWc)
}

// handleWc handles the wc command.
// All flags are boolean (no flag takes an argument).
// All positional arguments are file paths.
func handleWc(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	for _, tok := range args {
		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// All positionals are file paths
		result = append(result, "<path>")
	}

	result = append(result, redirects...)
	return result
}
