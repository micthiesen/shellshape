package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("chmod", handleChmod)
}

var chmodModeRE = regexp.MustCompile(`^[0-7]{3,4}$`)
var chmodSymbolicRE = regexp.MustCompile(`^[ugoa]*[=+-][rwxXstugo]*(,[ugoa]*[=+-][rwxXstugo]*)*$`)

// handleChmod handles the chmod command.
// The mode (octal or symbolic) is preserved verbatim because it defines
// the intent of the command. File arguments are collapsed to <path>.
func handleChmod(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	modeAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !modeAssigned {
				modeAssigned = true
			}
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the mode — keep verbatim.
		if !modeAssigned {
			if chmodModeRE.MatchString(tok) || chmodSymbolicRE.MatchString(tok) {
				result = append(result, tok)
			} else {
				// Unrecognized mode format, keep verbatim anyway.
				result = append(result, tok)
			}
			modeAssigned = true
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
