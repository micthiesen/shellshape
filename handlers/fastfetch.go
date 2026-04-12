package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("fastfetch", handleFastfetch)
}

// handleFastfetch handles the fastfetch system info tool.
// --logo keeps its argument verbatim (structural).
// Config flags collapse to <path>. Color/separator/structure collapse to <val>.
// Most flags are boolean.
func handleFastfetch(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	structuralFlags := map[string]bool{
		"--logo": true,
	}

	pathFlags := map[string]bool{
		"-c": true, "--config": true, "--load-config": true,
		"--data-raw": true,
	}

	valFlags := map[string]bool{
		"--color": true, "--separator": true,
		"--structure": true, "--key-type": true,
		"--format": true, "--title-fqdn": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
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

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
