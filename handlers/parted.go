package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("parted", handleParted)
}

// handleParted handles GNU parted.
// The first positional is the device path, collapsed to <device>.
// Flags -a/--align consume the next token as <val>.
// Inline command keywords (mklabel, mkpart, print, rm, etc.) are kept verbatim.
// The keyword "free" after "print" is also kept verbatim.
// All other positionals (command arguments: sizes, types, numbers, names) become <val>.
func handleParted(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a value.
	valFlags := map[string]bool{
		"-a": true, "--align": true,
	}

	// Known inline command keywords that stay verbatim.
	commands := map[string]bool{
		"mklabel":     true,
		"mkpart":      true,
		"print":       true,
		"rm":          true,
		"resizepart":  true,
		"set":         true,
		"toggle":      true,
		"name":        true,
		"unit":        true,
		"select":      true,
		"rescue":      true,
		"align-check": true,
		"disk_set":    true,
		"disk_toggle": true,
		"help":        true,
		"quit":        true,
	}

	// "free" is a keyword only when it follows "print".
	printModifiers := map[string]bool{
		"free": true,
		"list": true,
		"all":  true,
	}

	var result []string
	deviceAssigned := false
	lastCommand := ""

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			// A subshell in device position counts as device assigned.
			if !deviceAssigned {
				deviceAssigned = true
			}
			i++
			continue
		}

		// Flags that consume the next token.
		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Boolean flags.
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the device.
		if !deviceAssigned {
			result = append(result, "<device>")
			deviceAssigned = true
			i++
			continue
		}

		// Inline command keywords stay verbatim.
		if commands[tok] {
			result = append(result, tok)
			lastCommand = tok
			i++
			continue
		}

		// "free", "list", "all" after "print" stay verbatim.
		if lastCommand == "print" && printModifiers[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Everything else is a command argument, collapsed to <val>.
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
