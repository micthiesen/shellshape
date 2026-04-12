package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("dbus-monitor", handleDbusMonitor)
}

// handleDbusMonitor handles dbus-monitor (D-Bus message monitor).
// --system, --session, --profile, --monitor are boolean.
// --address → <val>.
// Positional filter expressions → <val>.
func handleDbusMonitor(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{"--address": true}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
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

		// Positional: filter expression → <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
