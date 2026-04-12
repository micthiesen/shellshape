package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("winetricks", handleWinetricks)
}

// handleWinetricks handles the winetricks Wine helper.
// Positionals are "verbs" (structural identifiers like dxvk, vcrun2019) → kept verbatim.
// --prefix consumes the next arg as <path>.
// Boolean flags: --unattended/-q, --force, --gui.
func handleWinetricks(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{"--prefix": true}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals are verbs (structural): keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
