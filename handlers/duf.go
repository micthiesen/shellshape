package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("duf", handleDuf)
}

// handleDuf handles the duf disk usage/free utility.
// Value flags (--hide, --only, --output) collapse their argument to <val>.
// Structural flags (--sort, --theme) keep their argument verbatim.
// Positionals are mount paths.
func handleDuf(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"--hide": true, "--only": true, "--output": true}, Placeholder: "<val>"},
	}

	structuralFlags := map[string]bool{"--sort": true, "--theme": true}

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

		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
