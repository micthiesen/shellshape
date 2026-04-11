package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"ln", "link"} {
		shellshape.Register(name, handleLn)
	}
}

// handleLn handles ln and link commands.
// All flags are boolean (-s, -f, -F, -L, -P, -h, -i, -n, -v, -w).
// All positionals are file paths (source files and target), classified via classifyToken.
// Subshells are preserved verbatim.
func handleLn(subcommand string, tokens []string) []string {
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
