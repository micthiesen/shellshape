package shellshape

import "regexp"

func init() {
	for _, name := range []string{"head", "tail"} {
		Register(name, handleTail)
	}
}

var tailNumericRE = regexp.MustCompile(`^[+-]?\d+[bckKMGTPE]?$`)

// handleTail handles the tail command.
// Numeric flags (-n, -c, -b) consume the next token as N.
// All positional arguments are file paths.
func handleTail(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "-c": true, "-b": true,
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
