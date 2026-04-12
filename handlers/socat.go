package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("socat", handleSocat)
}

func handleSocat(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-b": true, "-T": true,
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

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: address spec
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
