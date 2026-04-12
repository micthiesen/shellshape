package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("kscreen-doctor", handleKscreenDoctor)
}

func handleKscreenDoctor(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	for i := 0; i < len(args); i++ {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// Boolean flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// Positional: output configuration → <val>
		result = append(result, "<val>")
	}

	result = append(result, redirects...)
	return result
}
