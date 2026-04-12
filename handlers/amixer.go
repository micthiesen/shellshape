package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("amixer", handleAmixer, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleAmixer handles amixer subcommand arguments.
// Subcommands: scontrols, scontents, sset, sget, controls, contents, cset, cget.
// Control names and volume values collapse to <val>.
// -c <card> → N, -D <device> → <val>.
func handleAmixer(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-c": true, "--card": true,
	}
	valFlags := map[string]bool{
		"-D": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
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

		// All positionals are data (control names, values, specs)
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
