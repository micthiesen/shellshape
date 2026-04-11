package shellshape

func init() {
	for _, name := range []string{"echo", "printf"} {
		Register(name, handleEcho)
	}
}

// handleEcho handles echo and printf commands.
// Flags before the first positional are preserved. All positional string args
// collapse to a single <str>. Subshells are kept verbatim.
func handleEcho(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	hasPositional := false

	for _, tok := range args {
		if isFlagToken(tok) && !hasPositional {
			result = append(result, classifyToken(tok))
			continue
		}
		if isSubshellToken(tok) {
			if hasPositional {
				result = append(result, "<str>")
				hasPositional = false
			}
			result = append(result, tok)
			continue
		}
		hasPositional = true
	}
	if hasPositional {
		result = append(result, "<str>")
	}

	result = append(result, redirects...)
	return result
}
