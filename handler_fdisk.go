package shellshape

func init() {
	Register("fdisk", handleFdisk)
}

// handleFdisk handles the fdisk command.
// Numeric flags (-b, -C, -H, -S) consume the next token as N.
// The -s flag consumes the next token as <path> (partition device).
// All positionals (device paths) use classifyToken.
func handleFdisk(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-b": true, "-C": true, "-H": true, "-S": true,
	}
	pathFlags := map[string]bool{
		"-s": true,
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

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
