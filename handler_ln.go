package shellshape

// handleLn handles ln and link commands.
// All flags are boolean (-s, -f, -F, -L, -P, -h, -i, -n, -v, -w).
// All positionals are file paths (source files and target), classified via classifyToken.
// Subshells are preserved verbatim.
func handleLn(subcommand string, tokens []string) []string {
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
