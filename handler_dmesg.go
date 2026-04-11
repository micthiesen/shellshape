package shellshape

func init() {
	Register("dmesg", handleDmesg)
}

// handleDmesg handles the dmesg command (kernel ring buffer).
// Flags that consume arguments:
//
//	-f/--facility    → <facility>
//	-l/--level       → <level>
//	-n/--console-level → <level>
//	-F/--file, -K/--kmsg-file, -M, -N → <path>
//	-s/--buffer-size → N
//	--since/--until  → <time>
//	--time-format    → <fmt>
//
// All other flags are boolean. No positional arguments expected;
// any stray positionals get classifyToken.
func handleDmesg(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-F": true, "--file": true,
		"-K": true, "--kmsg-file": true,
		"-M": true, "-N": true,
	}
	levelFlags := map[string]bool{
		"-l": true, "--level": true,
		"-n": true, "--console-level": true,
	}
	facilityFlags := map[string]bool{
		"-f": true, "--facility": true,
	}
	sizeFlags := map[string]bool{
		"-s": true, "--buffer-size": true,
	}
	timeFlags := map[string]bool{
		"--since": true, "--until": true,
	}
	fmtFlags := map[string]bool{
		"--time-format": true,
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
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

		if levelFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<level>")
				}
				i++
			}
			continue
		}

		if facilityFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<facility>")
				}
				i++
			}
			continue
		}

		if sizeFlags[tok] {
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

		if timeFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<time>")
				}
				i++
			}
			continue
		}

		if fmtFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<fmt>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Unexpected positional
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
