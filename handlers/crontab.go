package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("crontab", handleCrontab)
}

// handleCrontab handles the crontab command.
// -u takes a user argument (collapsed to <user>).
// -l, -r, -e, -i are boolean flags.
// The sole positional argument is a file path to install.
func handleCrontab(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	userFlags := map[string]bool{"-u": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if userFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<user>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file to install
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
