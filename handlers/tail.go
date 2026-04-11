package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	for _, name := range []string{"head", "tail"} {
		shellshape.Register(name, handleTail)
	}
}

var tailNumericRE = regexp.MustCompile(`^[+-]?\d+[bckKMGTPE]?$`)

// handleTail handles the tail command.
// Numeric flags (-n, -c, -b) consume the next token as N.
// All positional arguments are file paths.
func handleTail(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "-c": true, "-b": true,
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
