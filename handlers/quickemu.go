package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("quickemu", handleQuickemu)
}

// handleQuickemu handles the quickemu QEMU wrapper.
// --vm <path>, --public-dir <path> collapse to <path>.
// Structural flags (--display, --keyboard, etc.) keep their argument verbatim.
// --snapshot takes an action (structural) and optionally a tag (<val>).
// Boolean flags (--fullscreen, --status-quo, etc.) are kept verbatim.
// Remaining positionals use ClassifyToken.
func handleQuickemu(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"--vm": true, "--public-dir": true,
	}

	structuralFlags := map[string]bool{
		"--display": true, "--keyboard": true, "--mouse": true,
		"--usb-controller": true, "--monitor": true, "--serial": true,
		"--viewer": true, "--sound-card": true, "--screen": true,
		"--accessible": true,
	}

	valFlags := map[string]bool{
		"--ssh-port": true, "--spice-port": true,
		"--monitor-telnet-host": true, "--monitor-telnet-port": true,
		"--serial-telnet-host": true, "--serial-telnet-port": true,
		"--keyboard_layout": true, "--monitor-cmd": true,
		"--extra_args": true,
	}

	// Snapshot actions that take a tag argument.
	snapshotWithTag := map[string]bool{
		"create": true, "apply": true, "delete": true,
	}

	categories := []shellshape.FlagCategory{
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

		// --snapshot action [tag]
		if tok == "--snapshot" {
			result = append(result, tok)
			i++
			if i < len(args) {
				action := args[i]
				result = append(result, action) // action is structural
				i++
				// If the action takes a tag, consume it.
				if snapshotWithTag[action] && i < len(args) {
					result = shellshape.EmitPositional(result, args[i], "<val>")
					i++
				}
			}
			continue
		}

		if structuralFlags[tok] {
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
