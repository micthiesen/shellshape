package shellshape

import "regexp"

func init() {
	Register("cut", handleCut)
}

var cutFusedRangeRE = regexp.MustCompile(`^-([bcf])(.+)$`)
var cutFusedDelimRE = regexp.MustCompile(`^-d(.+)$`)

// handleCut handles the cut command.
// -b, -c, -f take a range spec (e.g. 1-5,7) → <range>.
// -d takes a delimiter character → <delim>.
// -s, -n, -w are boolean flags.
// Positional arguments are file paths → classifyToken.
func handleCut(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	rangeFlags := map[string]bool{"-b": true, "-c": true, "-f": true}

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
		if tok == "-d" {
			result = append(result, "-d")
			i++
			if i < len(args) {
				result = append(result, "<delim>")
				i++
			}
			continue
		}

		// Fused delimiter: -d: → -d <delim>
		if m := cutFusedDelimRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-d", "<delim>")
			i++
			continue
		}

		// -b, -c, -f with separate argument
		if rangeFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<range>")
				i++
			}
			continue
		}

		// Fused range: -f1,3 → -f <range>
		if m := cutFusedRangeRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "<range>")
			i++
			continue
		}

		// Boolean flags and other flags
		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
