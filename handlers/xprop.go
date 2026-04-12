package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("xprop", handleXprop)
}

// handleXprop handles xprop (X window property displayer).
// -id, -name, -display → <val>
// -root, -spy, -notype, -frame, -remove, -len → boolean
// Positional property names are structural (kept verbatim).
func handleXprop(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-id": true, "-name": true, "-display": true,
		"-f": true,
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

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: property name, keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
