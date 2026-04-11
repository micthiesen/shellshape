package shellshape

// handleRm handles the rm (and unlink) command.
// All flags are boolean (no flag consumes an argument).
// After "--", all remaining tokens are treated as positionals.
// All positionals are file paths, classified via classifyToken.
// Subshells are preserved verbatim.
func handleRm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	endOfFlags := false

	for _, tok := range args {
		if tok == "--" {
			result = append(result, tok)
			endOfFlags = true
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if !endOfFlags && isFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// After "--", all tokens are positional file paths regardless of
		// whether they look like flags (e.g. "-weird-file").
		if endOfFlags {
			result = append(result, "<path>")
		} else {
			result = append(result, classifyToken(tok))
		}
	}

	result = append(result, redirects...)
	return result
}
