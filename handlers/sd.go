package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("sd", handleSd)
}

func handleSd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "--max-replacements": true,
	}

	var result []string
	positional := 0 // 0=find, 1=replace, 2+=files

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			positional++
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		switch positional {
		case 0:
			result = append(result, "<pattern>")
		case 1:
			result = append(result, "<val>")
		default:
			result = append(result, shellshape.ClassifyToken(tok))
		}
		positional++
		i++
	}

	result = append(result, redirects...)
	return result
}
