package shellshape

func init() {
	Register("mv", handleMv)
}

// handleMv handles the mv command.
// All flags are boolean except -t/--target-directory (consumes next arg as
// path) and -S/--suffix (consumes next arg as string).
// All positionals are file paths, collapsed to <path>.
func handleMv(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-t": true, "--target-directory": true,
	}
	strFlags := map[string]bool{
		"-S": true, "--suffix": true,
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

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if strFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a file path.
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
