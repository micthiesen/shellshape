package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"base64", "b64encode", "b64decode"} {
		shellshape.Register(name, handleBase64)
	}
}

// handleBase64 handles base64, b64encode, b64decode.
// -o/--output takes a path argument.
// -w/--wrap, -b/--break take a numeric argument.
// -d, -D, --decode, -i, -h, --help are boolean.
// Positionals are input file paths.
func handleBase64(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
	}
	numericFlags := map[string]bool{
		"-w": true, "--wrap": true,
		"-b": true, "--break": true,
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

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
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

		// Positional: input file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
