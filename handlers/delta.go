package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("delta", handleDelta)
}

// handleDelta handles the delta diff viewer.
// Numeric flags: --width, --tabs → N.
// Value flags: --paging → <val>.
// All other flags are boolean. Positionals are file paths.
func handleDelta(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"--width": true, "--tabs": true}, Placeholder: "N"},
		{Flags: map[string]bool{"--paging": true}, Placeholder: "<val>"},
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

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
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

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
