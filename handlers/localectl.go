package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("localectl", handleLocalectl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleLocalectl handles localectl subcommand arguments.
// All positional arguments (locale names, keymap names, x11 layout params) → <val>.
func handleLocalectl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-H": true, "--host": true,
		"-M": true, "--machine": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	// Handle global flag consumed as subcommand.
	if shellshape.IsFlagToken(subcommand) {
		if placeholder, ok := shellshape.MatchFlagCategory(subcommand, categories); ok {
			if i < len(args) && !shellshape.IsFlagToken(args[i]) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, placeholder)
				}
				i++
			}
		}
		if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
			result = append(result, args[i])
			i++
		}
	}

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

		// Positional → <val> (collapsed to single placeholder)
		if !hasPositional {
			result = append(result, "<val>")
			hasPositional = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
