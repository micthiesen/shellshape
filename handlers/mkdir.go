package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("mkdir", handleMkdir)
}

// handleMkdir handles the mkdir command.
// -m/--mode consumes the next token as <mode>.
// -p, -v are boolean flags.
// Bundled flags ending in 'm' (e.g. -pm) consume the next token as <mode>.
// All positionals use classifyToken. Subshells are preserved.
func handleMkdir(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	modeFlags := map[string]bool{"-m": true, "--mode": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if modeFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<mode>")
			continue
		}

		// Bundled short flags ending in 'm' (e.g. -pm, -pvm): the 'm' consumes next arg.
		if shellshape.IsFlagToken(tok) && !strings.HasPrefix(tok, "--") && strings.HasSuffix(tok, "m") {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<mode>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
