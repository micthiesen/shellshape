package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("htmlq", handleHtmlq)
}

func handleHtmlq(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-r": true, "--remove-nodes": true,
			"-B": true, "--base": true,
		}, Placeholder: "<val>"},
		{Flags: map[string]bool{
			"-f": true, "--filename": true,
		}, Placeholder: "<path>"},
	}

	structuralFlags := map[string]bool{
		"-a": true, "--attribute": true,
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

		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
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

		// Positional: CSS selector
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
