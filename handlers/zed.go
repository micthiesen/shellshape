package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("zed", handleZed)
}

// handleZed handles the Zed editor.
// All flags are boolean. Positional args are file paths (possibly path:line:col).
func handleZed(subcommand string, tokens []string) []string {
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
