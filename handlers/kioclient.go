package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"kioclient", "kioclient6"} {
		shellshape.Register(name, handleKioclient, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

// handleKioclient handles kioclient/kioclient6 subcommand arguments.
// Subcommands: exec, cat, copy, move, remove, stat, ls, mkdir, appmenu.
// Positional args after the subcommand are URLs or paths, handled by ClassifyToken.
func handleKioclient(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			continue
		}
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
