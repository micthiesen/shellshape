package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("compsize", handleCompsize)
}

// handleCompsize handles compsize (btrfs compression statistics).
//
// Simple command: all positional arguments are filesystem paths → <path>.
// Only common flag is -x (one filesystem, boolean).
func handleCompsize(subcommand string, tokens []string) []string {
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

		// All positionals are paths.
		result = append(result, "<path>")
	}

	result = append(result, redirects...)
	return result
}
