package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("tac", handleTac)
}

// handleTac handles the tac command.
// -b and -r are boolean flags.
// -s / --separator consumes the next token as <sep>.
// --separator=VALUE is rewritten to --separator=<sep>.
// All positionals are file paths, classified via classifyToken.
func handleTac(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	separatorFlags := map[string]bool{"-s": true, "--separator": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// --separator=VALUE
		if strings.HasPrefix(tok, "--separator=") {
			result = append(result, "--separator=<sep>")
			i++
			continue
		}

		if separatorFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<sep>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
