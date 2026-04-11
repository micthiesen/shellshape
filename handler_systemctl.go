package shellshape

func init() {
	Register("systemctl", handleSystemctl, HandlerOptions{HasSubcommands: true})
}

// handleSystemctl handles systemctl subcommand arguments.
// Since systemctl has HasSubcommands, the normalizer extracts the first token
// after the executable as the "subcommand". However, systemctl allows global
// flags before the subcommand (e.g. "systemctl --user start nginx"), so the
// normalizer may consume a flag as the subcommand. This handler detects that
// case and emits the real subcommand from the remaining tokens.
//
// Unit names (positional arguments) are collapsed to a single <unit>.
// Type and state flag values are kept verbatim (small finite sets).
// Property, host, and machine flag values are collapsed to <val>.
// Lines flag values (-n/--lines) are collapsed to N.
// Root flag values are collapsed to <path>.
func handleSystemctl(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose values are kept verbatim (small finite set of values).
	verbatimValueFlags := map[string]bool{
		"-t": true, "--type": true,
		"--state":  true,
		"--output": true,
	}

	// Flags whose values collapse to <val>.
	valFlags := map[string]bool{
		"-p": true, "--property": true,
		"-H": true, "--host": true,
		"-M": true, "--machine": true,
		"--kill-who":    true,
		"--signal":      true,
		"-s":            true,
		"--job-mode":    true,
		"--preset-mode": true,
	}

	// Flags whose values collapse to N.
	numericFlags := map[string]bool{
		"-n": true, "--lines": true,
	}

	// Flags whose values collapse to <path>.
	pathFlags := map[string]bool{
		"--root": true,
	}

	// Boolean flags (no argument consumed).
	booleanFlags := map[string]bool{
		"--user": true, "--system": true,
		"--failed": true,
		"--all":    true, "-a": true,
		"-l": true, "--full": true,
		"--no-pager": true, "--no-wall": true,
		"--no-block": true, "--no-reload": true,
		"--no-legend": true, "--no-ask-password": true,
		"-q": true, "--quiet": true,
		"--force": true, "-f": true,
		"--recursive":        true,
		"--value":            true,
		"--now":              true,
		"--runtime":          true,
		"--plain":            true,
		"--reverse":          true,
		"--wait":             true,
		"--show-transaction": true,
	}

	var result []string
	i := 0

	// When the normalizer consumed a global flag as the "subcommand", the
	// real subcommand is still in args. Consume any leading flags/values and
	// emit the real subcommand.
	if isFlagToken(subcommand) {
		// The flag token itself was already appended by the normalizer.
		// If it takes a value, the value is the first token in args.
		if valFlags[subcommand] && i < len(args) && !isFlagToken(args[i]) {
			if isSubshellToken(args[i]) {
				result = append(result, args[i])
			} else {
				result = append(result, "<val>")
			}
			i++
		} else if numericFlags[subcommand] && i < len(args) && !isFlagToken(args[i]) {
			if isSubshellToken(args[i]) {
				result = append(result, args[i])
			} else {
				result = append(result, "N")
			}
			i++
		} else if pathFlags[subcommand] && i < len(args) && !isFlagToken(args[i]) {
			if isSubshellToken(args[i]) {
				result = append(result, args[i])
			} else {
				result = append(result, "<path>")
			}
			i++
		} else if verbatimValueFlags[subcommand] && i < len(args) && !isFlagToken(args[i]) {
			result = append(result, args[i])
			i++
		}
		// Now emit the real subcommand (next non-flag positional).
		if i < len(args) && !isFlagToken(args[i]) && !isSubshellToken(args[i]) {
			result = append(result, args[i])
			i++
		}
	}

	hasUnit := false

	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if booleanFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if verbatimValueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		// Other flags (unknown flags): keep verbatim.
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: unit name, collapse all to a single <unit>.
		if !hasUnit {
			result = append(result, "<unit>")
			hasUnit = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
