package shellshape

func init() {
	Register("df", handleDf)
}

// handleDf handles the df (disk free) command.
// -T/--type and -t (legacy) consume a filesystem type value (<type>).
// -x/--exclude-type consumes a filesystem type value (<type>).
// Combined short flags ending in T/t/x (e.g. -lT) consume the next arg as <type>.
// All other flags are boolean. Positionals are file/filesystem paths via classifyToken.
func handleDf(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	typeFlags := map[string]bool{
		"-T": true, "--type": true,
		"-t": true,
		"-x": true, "--exclude-type": true,
	}

	// Short flag letters that consume a type argument.
	typeLetters := map[byte]bool{'T': true, 't': true, 'x': true}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if typeFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<type>")
			continue
		}

		// Combined short flags like -lT: if last char is a type-consuming letter,
		// the next token is the type value.
		if len(tok) > 2 && tok[0] == '-' && tok[1] != '-' && typeLetters[tok[len(tok)-1]] {
			result, i = consumeFlagArg(tok, args, i, result, "<type>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: classify as path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
