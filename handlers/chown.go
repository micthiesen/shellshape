package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("chown", handleChown)
}

// handleChown handles the chown command.
// The owner[:group] spec is preserved verbatim because it defines
// the intent of the command. File arguments are collapsed to <path>.
func handleChown(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	ownerAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !ownerAssigned {
				ownerAssigned = true
			}
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the owner[:group] — keep verbatim.
		if !ownerAssigned {
			result = append(result, tok)
			ownerAssigned = true
			i++
			continue
		}

		// Remaining positionals are file paths.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
