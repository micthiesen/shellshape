package shellshape

// handleLs handles the ls command.
// -D is the only flag that consumes the next token (date format string).
// All other flags are boolean. All positionals are paths.
func handleLs(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	formatFlags := map[string]bool{"-D": true}

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

		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
