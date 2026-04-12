package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("ffplay", handleFfplay)
}

// handleFfplay handles the ffplay media player.
// -ss, -t collapse to <val> (time values).
// -volume, -x, -y, -loop collapse to N (numeric).
// -v (loglevel) keeps value verbatim (structural).
// -autoexit, -fs, -an, -vn are boolean.
// Positionals are classified via ClassifyToken (paths or URLs).
func handleFfplay(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-ss": true, "-t": true}, Placeholder: "<val>"},
		{Flags: map[string]bool{"-volume": true, "-x": true, "-y": true, "-loop": true}, Placeholder: "N"},
	}

	// Flags whose next arg is kept verbatim
	verbatimFlags := map[string]bool{
		"-v": true,
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

		if verbatimFlags[tok] {
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

		// Positional: classify (path or URL)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
