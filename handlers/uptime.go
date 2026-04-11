package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("uptime", handleUptime)
}

// handleUptime handles the uptime command.
// All flags are boolean (-p, -s, -V, -h and long forms). Any unexpected
// positional arguments are classified generically. Subshells are preserved.
func handleUptime(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
