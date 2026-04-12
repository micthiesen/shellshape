package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("wpctl", handleWpctl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleWpctl handles wpctl subcommand arguments.
// Subcommands: status, inspect, set-volume, set-mute, set-default, set-profile, get-volume.
// Node IDs, volume values, and profile names all collapse to <val>.
func handleWpctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// All positionals are data
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
