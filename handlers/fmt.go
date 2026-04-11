package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("fmt", handleFmt)
}

// handleFmt handles the fmt text formatting command.
// Value Flags: -d (sentence-ending chars, string), -l (tab width, numeric),
// -t (input tab stops, numeric), -w (line width, numeric).
// Boolean Flags: -c, -m, -n, -p, -s.
// All positionals are input file paths.
func handleFmt(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-w": true, "-l": true, "-t": true,
	}
	stringFlags := map[string]bool{
		"-d": true,
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

		if stringFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: input file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
