package shellshape

import "regexp"

func init() {
	Register("paste", handlePaste)
}

var pasteFusedDelimRE = regexp.MustCompile(`^-d(.+)$`)

// handlePaste handles the paste command.
// -d takes a delimiter string → <delim>.
// -s/--serial is a boolean flag.
// All positionals are file paths (- means stdin, kept verbatim).
func handlePaste(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -d with separate argument
		if tok == "-d" || tok == "--delimiters" {
			result, i = consumeFlagArg(tok, args, i, result, "<delim>")
			continue
		}

		// Fused delimiter: -d: → -d <delim>
		if m := pasteFusedDelimRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-d", "<delim>")
			i++
			continue
		}

		// Boolean flags and other flags
		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: file path (- stays verbatim via classifyToken)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
