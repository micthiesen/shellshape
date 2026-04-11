package shellshape

func init() {
	Register("who", handleWho)
}

// handleWho handles the who command.
// All flags are boolean and preserved verbatim.
// "am i" / "am I" is a special idiom preserved as structural.
// The optional file positional is classified generically.
func handleWho(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// "am i" / "am I" idiom: preserve verbatim.
		if (tok == "am" || tok == "AM") && i+1 < len(args) {
			next := args[i+1]
			if next == "i" || next == "I" {
				result = append(result, tok, next)
				i += 2
				continue
			}
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: optional file argument.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
