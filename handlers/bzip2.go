package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"bzip2", "bunzip2", "bzcat"} {
		shellshape.Register(name, handleBzip2)
	}
}

// handleBzip2 handles bzip2, bunzip2, bzcat.
// All flags are boolean (-d, -z, -k, -f, -v, -t, -c).
// Compression levels -1 through -9 are preserved verbatim (structural).
// Positionals are file paths.
func handleBzip2(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Compression levels: -1 through -9
		if len(tok) == 2 && tok[0] == '-' && tok[1] >= '1' && tok[1] <= '9' {
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
