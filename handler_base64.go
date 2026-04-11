package shellshape

// handleBase64 handles base64, b64encode, b64decode.
// -o/--output takes a path argument.
// -w/--wrap, -b/--break take a numeric argument.
// -d, -D, --decode, -i, -h, --help are boolean.
// Positionals are input file paths.
func handleBase64(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
	}
	numericFlags := map[string]bool{
		"-w": true, "--wrap": true,
		"-b": true, "--break": true,
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
				result = append(result, "<path>")
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

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: input file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
