package shellshape

func init() {
	Register("whois", handleWhois)
}

// handleWhois handles the whois command.
// -h (host) consumes the next token as <host>.
// -p (port) consumes the next token as N.
// -c (TLD) consumes the next token as <val>.
// All other single-letter flags are boolean.
// All positional arguments collapse to a single <name>.
func handleWhois(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	hostFlags := map[string]bool{"-h": true}
	numericFlags := map[string]bool{"-p": true}
	valueFlags := map[string]bool{"-c": true}

	var result []string
	nameAdded := false
	pastDashes := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if tok == "--" && !pastDashes {
			result = append(result, tok)
			pastDashes = true
			i++
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if !pastDashes && hostFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<host>")
			continue
		}

		if !pastDashes && numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if !pastDashes && valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if !pastDashes && isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: collapse all to single <name>
		if !nameAdded {
			result = append(result, "<name>")
			nameAdded = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
