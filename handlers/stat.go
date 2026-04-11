package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("stat", handleStat)
}

// handleStat handles the stat command.
// Format flags (-f, -t, -c, --format, --printf) consume the next arg as <fmt>.
// Boolean flags (-L, -x, -F, -l, -n, -q, -r, -s) are kept verbatim.
// All positionals are file paths and use classifyToken.
func handleStat(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	formatFlags := map[string]bool{
		"-f": true, "-t": true, "-c": true,
		"--format": true, "--printf": true,
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if formatFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<fmt>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
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
