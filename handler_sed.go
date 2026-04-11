package shellshape

func init() {
	Register("sed", handleSed)
}

// handleSed handles the sed command.
// Expression flags (-e, --expression) make the next arg <sed-expr>.
// File flags (-f, --file) make the next arg <path>.
// First positional (if no script assigned) becomes <sed-expr>.
// Remaining positionals use classifyToken.
func handleSed(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	sedExprFlags := map[string]bool{"-e": true, "--expression": true}
	sedFileFlags := map[string]bool{"-f": true, "--file": true}

	var result []string
	scriptAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			scriptAssigned = true
			i++
			continue
		}

		if sedExprFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<sed-expr>")
				scriptAssigned = true
				i++
			}
			continue
		}

		if sedFileFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<script>")
				scriptAssigned = true
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !scriptAssigned {
			result = append(result, "<sed-expr>")
			scriptAssigned = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
