package shellshape

func init() {
	Register("crontab", handleCrontab)
}

// handleCrontab handles the crontab command.
// -u takes a user argument (collapsed to <user>).
// -l, -r, -e, -i are boolean flags.
// The sole positional argument is a file path to install.
func handleCrontab(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	userFlags := map[string]bool{"-u": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if userFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<user>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file to install
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
