package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("protontricks", handleProtontricks)
}

// handleProtontricks handles the protontricks Proton helper.
// First positional is an app ID (number) → N.
// Remaining positionals are winetricks verbs (structural) → kept verbatim.
// -s/--search and -c consume the next arg as <val>.
// Boolean flags: --no-runtime, --gui, -v/--verbose.
func handleProtontricks(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-s": true, "--search": true,
		"-c": true,
	}

	var result []string
	appIDFound := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			appIDFound = true
			i++
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional: app ID (number) → N
		if !appIDFound {
			if shellshape.NumberRE.MatchString(tok) {
				result = append(result, "N")
			} else {
				result = append(result, tok)
			}
			appIDFound = true
			i++
			continue
		}

		// Remaining positionals are verbs (structural): keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
