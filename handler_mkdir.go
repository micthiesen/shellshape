package shellshape

import "strings"

func init() {
	Register("mkdir", handleMkdir)
}

// handleMkdir handles the mkdir command.
// -m/--mode consumes the next token as <mode>.
// -p, -v are boolean flags.
// Bundled flags ending in 'm' (e.g. -pm) consume the next token as <mode>.
// All positionals use classifyToken. Subshells are preserved.
func handleMkdir(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	modeFlags := map[string]bool{"-m": true, "--mode": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if modeFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<mode>")
				i++
			}
			continue
		}

		// Bundled short flags ending in 'm' (e.g. -pm, -pvm): the 'm' consumes next arg.
		if isFlagToken(tok) && !strings.HasPrefix(tok, "--") && strings.HasSuffix(tok, "m") {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<mode>")
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
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
