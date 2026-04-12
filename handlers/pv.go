package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("pv", handlePv)
}

func handlePv(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-s": true, "--size": true,
			"-L": true, "--rate-limit": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{
			"-N": true, "--name": true,
		}, Placeholder: "<val>"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a file path for pv
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
