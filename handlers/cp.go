package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("cp", handleCp)
}

// handleCp handles the cp command.
// All flags are boolean (no flag consumes an argument).
// All positionals are file paths, classified via classifyToken.
// Subshells are preserved verbatim.
func handleCp(subcommand string, tokens []string) []string {
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
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
