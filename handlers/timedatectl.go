package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("timedatectl", handleTimedatectl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleTimedatectl handles timedatectl subcommand arguments.
// All positional arguments (timezones, time strings, NTP booleans) → <val>.
func handleTimedatectl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-H": true, "--host": true,
		"-M": true, "--machine": true,
		"-p": true, "--property": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	result, _, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, categories, nil,
	)

	hasPositional := false

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

		// Positional → <val>
		if !hasPositional {
			result = append(result, "<val>")
			hasPositional = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
