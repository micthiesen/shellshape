package shellshape

func init() {
	Register("fsck", handleFsck)
}

// handleFsck handles the fsck (filesystem check) command.
// -t consumes the next token as <type> (filesystem type).
// All other flags are boolean. Positionals are device/filesystem paths
// collapsed via classifyToken.
func handleFsck(subcommand string, tokens []string) []string {
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

		// -t/--type consumes the next token as a filesystem type.
		if tok == "-t" || tok == "--type" {
			result, i = consumeFlagArg(tok, args, i, result, "<type>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device or filesystem path.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
