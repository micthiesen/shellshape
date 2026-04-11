package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("tee", handleTee)
}

// handleTee handles the tee command.
// All flags are boolean (-a, -i). All positionals are output file paths
// and collapse to <path>. Subshells are kept verbatim.
func handleTee(subcommand string, tokens []string) []string {
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

		// All positionals are file paths.
		result = append(result, "<path>")
	}

	result = append(result, redirects...)
	return result
}
