package shellshape

func init() {
	Register("touch", handleTouch)
}

// handleTouch handles the touch command.
// -A takes a time adjustment value, -d takes a date string, -t takes a timestamp,
// -r takes a reference file path. Boolean flags: -a, -c, -h, -m.
// All positional arguments are file paths.
func handleTouch(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags where the next token is a value (timestamp, date, adjustment).
	valueFlags := map[string]bool{
		"-A": true, "-t": true, "-d": true,
	}
	// Flags where the next token is a file path.
	pathFlags := map[string]bool{
		"-r": true,
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

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
