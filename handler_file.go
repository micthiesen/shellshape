package shellshape

// handleFile handles the file command (determine file type).
// Flags -m/-M/--magic-file, -f/--files-from consume a path argument.
// Flags -F/--separator, -e/--exclude, -P consume a value argument.
// All positional arguments are file paths.
func handleFile(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-m": true, "--magic-file": true,
		"-M": true,
		"-f": true, "--files-from": true,
	}
	valueFlags := map[string]bool{
		"-F": true, "--separator": true,
		"-e": true, "--exclude": true,
		"-P": true,
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

		if valueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
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
