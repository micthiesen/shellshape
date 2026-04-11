package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"cat", "rev"} {
		shellshape.Register(name, handleCat)
	}
}

// handleCat handles the cat command.
// All flags are boolean (no flag consumes an argument).
// All positionals are file paths, classified via classifyToken.
// Subshells are preserved verbatim.
func handleCat(subcommand string, tokens []string) []string {
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
