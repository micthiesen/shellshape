package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("efibootmgr", handleEfibootmgr)
}

// handleEfibootmgr handles efibootmgr flags.
// Numeric: -b, -n, -t, -p (boot number, next boot, timeout, partition).
// Path: -d (disk), -l (loader path).
// Val: -o (boot order), -L (label).
// Boolean: -B, --create, -v.
func handleEfibootmgr(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-b": true, "-n": true, "-t": true, "-p": true}, Placeholder: "N"},
		{Flags: map[string]bool{"-d": true, "-l": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"-o": true, "-L": true}, Placeholder: "<val>"},
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

		// No expected positionals, but classify if present
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
