package shellshape

func init() {
	Register("tee", handleTee)
}

// handleTee handles the tee command.
// All flags are boolean (-a, -i). All positionals are output file paths
// and collapse to <path>. Subshells are kept verbatim.
func handleTee(subcommand string, tokens []string) []string {
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

		// All positionals are file paths.
		result = append(result, "<path>")
	}

	result = append(result, redirects...)
	return result
}
