package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("envsubst", handleEnvsubst)
}

// handleEnvsubst handles the envsubst command.
// The only flag is -v/--variables (boolean). An optional positional
// SHELL-FORMAT string collapses to <val>. Subshells are kept verbatim.
func handleEnvsubst(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}
		// Positional: SHELL-FORMAT string → <val>
		result = append(result, "<val>")
	}

	result = append(result, redirects...)
	return result
}
