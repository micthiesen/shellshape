package shellshape

func init() {
	for _, name := range []string{"traceroute", "tracepath"} {
		Register(name, handleTraceroute)
	}
}

// handleTraceroute handles traceroute and tracepath.
// Numeric flags (-m, -q, -w, -p, -f, -M, -t, -z) collapse next arg to N.
// Value flags (-s, -i, -g, -A, -P) collapse next arg to <val>.
// First positional becomes <host>, optional second positional becomes N (packet size).
func handleTraceroute(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-m": true, "-q": true, "-w": true, "-p": true,
		"-f": true, "-M": true, "-t": true, "-z": true,
	}

	valueFlags := map[string]bool{
		"-s": true, "-i": true, "-g": true, "-A": true, "-P": true,
	}

	// These don't pass isFlagToken (start with digit), so handle explicitly.
	booleanDigitFlags := map[string]bool{
		"-4": true, "-6": true,
	}

	var result []string
	hostAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if booleanDigitFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals: first is host, second is packet size (numeric)
		if !hostAssigned {
			result = append(result, "<host>")
			hostAssigned = true
		} else {
			result = append(result, "N")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
