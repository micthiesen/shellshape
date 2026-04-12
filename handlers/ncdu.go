package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("ncdu", handleNcdu)
}

// handleNcdu handles the ncdu disk usage analyzer.
// Path flags (-o, -f) collapse their argument to <path>.
// Exclude flags (--exclude) collapse their argument to <pattern>.
// Structural flags (--color) keep their argument verbatim.
// Positionals are paths.
func handleNcdu(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-o": true, "-f": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"--exclude": true}, Placeholder: "<pattern>"},
	}

	structuralFlags := map[string]bool{"--color": true}

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
