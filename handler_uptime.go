package shellshape

func init() {
	Register("uptime", handleUptime)
}

// handleUptime handles the uptime command.
// All flags are boolean (-p, -s, -V, -h and long forms). Any unexpected
// positional arguments are classified generically. Subshells are preserved.
func handleUptime(subcommand string, tokens []string) []string {
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
		result = append(result, classifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
