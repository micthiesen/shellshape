package shellshape

// handleCat handles the cat command.
// All flags are boolean (no flag consumes an argument).
// All positionals are file paths, classified via classifyToken.
// Subshells are preserved verbatim.
func handleCat(subcommand string, tokens []string) []string {
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
		result = append(result, classifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
