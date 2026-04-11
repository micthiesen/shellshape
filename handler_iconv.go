package shellshape

func init() {
	Register("iconv", handleIconv)
}

// handleIconv handles iconv (character encoding converter).
// -f/--from-code and -t/--to-code consume the next token as <encoding>.
// -o/--output consumes the next token as <path>.
// Remaining positionals are input files (classifyToken).
func handleIconv(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	encodingFlags := map[string]bool{
		"-f": true, "--from-code": true,
		"-t": true, "--to-code": true,
	}
	pathFlags := map[string]bool{
		"-o": true, "--output": true,
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

		if encodingFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<encoding>")
			continue
		}

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: input file
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
