package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"echo", "printf"} {
		shellshape.Register(name, handleEcho)
	}
}

// handleEcho handles echo and printf commands.
// Flags before the first positional are preserved. All positional string args
// collapse to a single <str>. Subshells are kept verbatim.
func handleEcho(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	hasPositional := false

	for _, tok := range args {
		if shellshape.IsFlagToken(tok) && !hasPositional {
			result = append(result, shellshape.ClassifyToken(tok))
			continue
		}
		if shellshape.IsSubshellToken(tok) {
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
