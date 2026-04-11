package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("paste", handlePaste)
}

var pasteFusedDelimRE = regexp.MustCompile(`^-d(.+)$`)

// handlePaste handles the paste command.
// -d takes a delimiter string → <delim>.
// -s/--serial is a boolean flag.
// All positionals are file paths (- means stdin, kept verbatim).
func handlePaste(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -d with separate argument
		if tok == "-d" || tok == "--delimiters" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<delim>")
			continue
		}

		// Fused delimiter: -d: → -d <delim>
		if m := pasteFusedDelimRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-d", "<delim>")
			i++
			continue
		}

		// Boolean flags and other flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: file path (- stays verbatim via classifyToken)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
