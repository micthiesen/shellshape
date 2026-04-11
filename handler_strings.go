package shellshape

func init() {
	Register("strings", handleStrings)
}

// handleStrings handles the strings command.
// Flags -n/--bytes take a numeric argument (-> N).
// Flags -t/--radix take a value argument (o/d/x), preserved verbatim.
// Flags -e/--encoding take a value argument (s/S/b/l/B/L), preserved verbatim.
// Flag -arch takes a value argument (-> <val>).
// Boolean flags: -a, -f, -o, -.
// All positionals are file paths (-> <path>).
func handleStrings(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "--bytes": true,
	}
	// Flags whose argument is a small enum; preserve verbatim.
	enumFlags := map[string]bool{
		"-t": true, "--radix": true,
		"-e": true, "--encoding": true,
	}
	// Flags whose argument is opaque data.
	valueFlags := map[string]bool{
		"-arch": true,
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

		if numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if enumFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if tok == "--" {
			result = append(result, tok)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a file path.
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
