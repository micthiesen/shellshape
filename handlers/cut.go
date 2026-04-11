package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("cut", handleCut)
}

var cutFusedRangeRE = regexp.MustCompile(`^-([bcf])(.+)$`)
var cutFusedDelimRE = regexp.MustCompile(`^-d(.+)$`)

// handleCut handles the cut command.
// -b, -c, -f take a range spec (e.g. 1-5,7) → <range>.
// -d takes a delimiter character → <delim>.
// -s, -n, -w are boolean flags.
// Positional arguments are file paths → classifyToken.
func handleCut(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	rangeFlags := map[string]bool{"-b": true, "-c": true, "-f": true}

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
		if tok == "-d" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<delim>")
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
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<range>")
			continue
		}

		// Fused range: -f1,3 → -f <range>
		if m := cutFusedRangeRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "<range>")
			i++
			continue
		}

		// Boolean flags and other flags
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
