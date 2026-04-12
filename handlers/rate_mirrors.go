package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("rate-mirrors", handleRateMirrors, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleRateMirrors handles rate-mirrors (mirror speed tester).
// Subcommands: arch, cachyos, chaotic-aur, endeavouros, manjaro.
// --protocol keeps its argument verbatim (structural: https, http).
// --country → <val>, --save → <path>.
// --max-delay, --fetch-mirrors-timeout, --concurrency, --per-mirror-timeout, --top → N.
func handleRateMirrors(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"--country": true, "--url": true,
			"--completion": true, "--entry-country": true,
		}, Placeholder: "<val>"},
		{Flags: map[string]bool{
			"--save": true,
		}, Placeholder: "<path>"},
		{Flags: map[string]bool{
			"--max-delay": true, "--fetch-mirrors-timeout": true,
			"--concurrency": true, "--per-mirror-timeout": true,
			"--top": true, "--eps": true,
		}, Placeholder: "N"},
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

		// --protocol keeps its argument verbatim (structural)
		if tok == "--protocol" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
