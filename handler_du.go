package shellshape

// handleDu handles the du (disk usage) command.
// -d/--max-depth consume a numeric argument (N).
// -B/--block-size consume a numeric argument (N).
// -t/--threshold consumes the next arg as <threshold> (can be e.g. "100M").
// -I/--exclude consume the next arg as <pattern>.
// All other flags are boolean. Positionals are file paths via classifyToken.
func handleDu(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-d": true, "--max-depth": true,
		"-B": true, "--block-size": true,
	}
	thresholdFlags := map[string]bool{
		"-t": true, "--threshold": true,
	}
	patternFlags := map[string]bool{
		"-I": true, "--exclude": true,
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
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if thresholdFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<threshold>")
				i++
			}
			continue
		}

		if patternFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<pattern>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: classify as path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
