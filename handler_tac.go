package shellshape

import "strings"

func init() {
	Register("tac", handleTac)
}

// handleTac handles the tac command.
// -b and -r are boolean flags.
// -s / --separator consumes the next token as <sep>.
// --separator=VALUE is rewritten to --separator=<sep>.
// All positionals are file paths, classified via classifyToken.
func handleTac(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	separatorFlags := map[string]bool{"-s": true, "--separator": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// --separator=VALUE
		if strings.HasPrefix(tok, "--separator=") {
			result = append(result, "--separator=<sep>")
			i++
			continue
		}

		if separatorFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<sep>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
