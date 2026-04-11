package shellshape

// handleTr handles the tr (translate characters) command.
// All flags (-C, -c, -d, -s, -u) are boolean and preserved verbatim.
// Positional arguments (string1 and string2) are character set specs
// and collapse to <set>. Subshells are kept verbatim.
func handleTr(subcommand string, tokens []string) []string {
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

		// Positional: character set argument
		result = append(result, "<set>")
	}

	result = append(result, redirects...)
	return result
}
