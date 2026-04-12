package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"kcmshell5", "kcmshell6"} {
		shellshape.Register(name, handleKcmshell)
	}
}

// handleKcmshell handles kcmshell5 and kcmshell6.
// Module names are structural (kept verbatim).
// --args consumes the next token as <val>.
// --list is boolean.
func handleKcmshell(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{"--args": true}

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

		// Positional: module name, keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
