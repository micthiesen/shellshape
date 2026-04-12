package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("xdg-open", handleXdgOpen)
}

// handleXdgOpen handles xdg-open commands.
// xdg-open takes a single positional argument (URL or path) and has no
// meaningful flags. We use ClassifyToken to collapse URLs and paths.
func handleXdgOpen(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
