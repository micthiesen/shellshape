package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("lspci", handleLspci)
}

// handleLspci handles lspci (PCI device listing).
// Most flags are boolean. A few consume arguments: slot/device filters → <val>,
// file paths → <path>.
func handleLspci(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-s": true, "-d": true, "-A": true,
		}, Placeholder: "<val>"},
		{Flags: map[string]bool{
			"-i": true, "-p": true, "-F": true,
		}, Placeholder: "<path>"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// No expected positionals; classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
