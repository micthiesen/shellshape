package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("picocom", handlePicocom)
}

func handlePicocom(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-b": true, "--baud": true,
		"--databits": true, "--stopbits": true,
	}
	// Structural flags: value preserved verbatim (flow type, parity type)
	structuralFlags := map[string]bool{
		"--flow": true, "--parity": true,
	}
	valFlags := map[string]bool{
		"--imap": true, "--omap": true, "--emap": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: structuralFlags, Placeholder: ""},
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

		// Check flag categories
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			if placeholder == "" {
				// Structural: preserve flag and value verbatim
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
			} else {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			}
			continue
		}

		// Fused long flags
		if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, norm)
			i++
			continue
		}

		// Boolean flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
