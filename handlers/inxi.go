package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("inxi", handleInxi)
}

// handleInxi handles inxi system information tool.
// Most flags are boolean. A few take numeric or path arguments.
func handleInxi(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-c": true, "-y": true,
		"--width": true, "--indent": true,
	}

	pathFlags := map[string]bool{
		"--output-file": true, "--log-file": true,
	}

	valFlags := map[string]bool{
		"--output": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
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
