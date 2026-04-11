package shellshape

func init() {
	Register("uname", handleUname)
}

// handleUname handles the uname command.
// All flags are boolean (no flag consumes an argument).
// Unexpected positionals are collapsed to <str>; subshells are preserved.
func handleUname(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		if isFlagToken(tok) {
			result = append(result, tok)
			continue
		}
		// Unexpected positional → collapse to <str>
		result = append(result, "<str>")
	}

	result = append(result, redirects...)
	return result
}
