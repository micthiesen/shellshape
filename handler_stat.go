package shellshape

// handleStat handles the stat command.
// Format flags (-f, -t, -c, --format, --printf) consume the next arg as <fmt>.
// Boolean flags (-L, -x, -F, -l, -n, -q, -r, -s) are kept verbatim.
// All positionals are file paths and use classifyToken.
func handleStat(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	formatFlags := map[string]bool{
		"-f": true, "-t": true, "-c": true,
		"--format": true, "--printf": true,
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

		if formatFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<fmt>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
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
