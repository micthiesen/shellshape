package shellshape

func init() {
	Register("host", handleHost)
}

// handleHost handles the host DNS lookup utility.
// -t (type) and -c (class) consume the next token verbatim (structural).
// -R, -W, -N consume the next token as N (numeric).
// -m consumes the next token verbatim (debug flag keyword).
// First positional becomes <name>, second becomes <server>.
// Boolean flags: -a, -C, -d, -i, -l, -n, -r, -s, -T, -v, -V, -w, -4, -6.
func handleHost(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	verbatimFlags := map[string]bool{"-t": true, "-c": true, "-m": true}
	numericFlags := map[string]bool{"-R": true, "-W": true, "-N": true}

	var result []string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if isFlagToken(tok) || tok == "-4" || tok == "-6" {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first is <name>, second is <server>
		positionalCount++
		if positionalCount == 1 {
			result = append(result, "<name>")
		} else {
			result = append(result, "<server>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
