package shellshape

func init() {
	Register("free", handleFree)
}

// handleFree handles the free (memory usage) command.
// -s/--seconds and -c/--count consume the next token as a numeric value (N).
// Combined short flags ending in s or c (e.g. -hs) consume the next token.
// All other flags are boolean. free takes no positional arguments.
func handleFree(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	consumingFlags := map[string]bool{
		"-s": true, "--seconds": true,
		"-c": true, "--count": true,
	}

	consumingLetters := map[byte]bool{'s': true, 'c': true}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if consumingFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		// Combined short flags like -hs: if last char is a consuming letter,
		// the next token is the value.
		if len(tok) > 2 && tok[0] == '-' && tok[1] != '-' && consumingLetters[tok[len(tok)-1]] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
