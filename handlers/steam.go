package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("steam", handleSteam)
}

// handleSteam handles the Steam client CLI.
// -applaunch <appid> collapses appid to N (numeric).
// -silent, -shutdown, -bigpicture, -console are boolean.
// Positionals are classified via ClassifyToken (steam:// URLs become <steam-uri>).
// Extra arguments after -applaunch N are game launch options, kept as flags.
func handleSteam(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-applaunch": true}, Placeholder: "N"},
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

		// Positional: classify (steam:// URLs, etc.)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
