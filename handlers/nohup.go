package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("nohup", handleNohup)
}

// handleNohup handles the nohup command.
// nohup takes no flags of its own. The first positional is the wrapped command
// name (preserved verbatim), and all remaining tokens are that command's
// arguments, classified generically via classifyToken. Subshells are kept
// verbatim.
func handleNohup(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	commandFound := false

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if !commandFound {
			// First positional is the wrapped command name.
			result = append(result, shellshape.ClassifyToken(tok))
			commandFound = true
			continue
		}

		// Remaining tokens are the wrapped command's arguments.
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
