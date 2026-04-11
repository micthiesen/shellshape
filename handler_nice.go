package shellshape

func init() {
	Register("nice", handleNice)
}

// handleNice handles the nice command.
// nice runs a utility with an altered scheduling priority.
// The -n flag consumes the next token as N (the niceness increment).
// The first non-flag positional is the utility name (preserved verbatim).
// Remaining tokens after the utility are classified generically via classifyToken.
func handleNice(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	utilityFound := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Once we've found the utility, everything else is utility args.
		if utilityFound {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// -n flag: consumes the next token as niceness value
		if tok == "-n" || tok == "--adjustment" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		// Other flags (e.g. --help, --version)
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First non-flag positional is the utility name: keep verbatim.
		result = append(result, tok)
		utilityFound = true
		i++
	}

	result = append(result, redirects...)
	return result
}
