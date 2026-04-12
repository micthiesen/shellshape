package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("killall", handleKillall)
}

// handleKillall handles the killall command.
// Signal flags (-9, -HUP, -SIGTERM, -s/--signal NAME) are kept verbatim.
// User flag (-u/--user) collapses to <val>.
// Boolean flags: -e/--exact, -i/--interactive, -g/--process-group, -r/--regexp,
// -w/--wait, -q/--quiet, -I (case insensitive).
// Positional process names collapse to <pattern>.
func handleKillall(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is structural (signal name, kept verbatim)
	signalFlags := map[string]bool{
		"-s": true, "--signal": true,
	}

	// Flags whose next token collapses to <val>
	valFlags := map[string]bool{
		"-u": true, "--user": true,
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

		if signalFlags[tok] {
			// Keep signal name verbatim
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Signal as flag: -9, -15, -HUP, -KILL, -SIGTERM, etc.
		if shellshape.IsFlagToken(tok) || isNumericSignal(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: process name
		result = append(result, "<pattern>")
		i++
	}

	result = append(result, redirects...)
	return result
}
