package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("ffprobe", handleFfprobe)
}

// handleFfprobe handles ffprobe commands.
// -v (loglevel), -of/-print_format (output format), -show_entries keep their values verbatim (structural).
// -select_streams collapses to <val>.
// -i (input) collapses to <path>.
// -show_format, -show_streams are boolean.
// Positionals are classified via ClassifyToken (paths become <path>, URLs become <scheme-uri>).
func handleFfprobe(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next arg is kept verbatim (structural)
	verbatimFlags := map[string]bool{
		"-v": true, "-of": true, "-print_format": true,
		"-show_entries": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-i": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"-select_streams": true}, Placeholder: "<val>"},
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
