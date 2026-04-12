package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("mangohud", handleMangohud)
}

// handleMangohud handles the mangohud performance overlay wrapper.
// Boolean flags: --dlsym.
// First non-flag positional is the wrapped command name (kept verbatim).
// Remaining tokens are the wrapped command's arguments, classified generically.
func handleMangohud(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	commandFound := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			commandFound = true
			i++
			continue
		}

		// Once we've found the command, everything else is its args.
		if commandFound {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Boolean flags before the command
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First non-flag positional is the wrapped command: keep verbatim.
		result = append(result, tok)
		commandFound = true
		i++
	}

	result = append(result, redirects...)
	return result
}
