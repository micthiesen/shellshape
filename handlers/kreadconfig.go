package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"kreadconfig6", "kwriteconfig6", "kreadconfig5", "kwriteconfig5"} {
		shellshape.Register(name, handleKreadconfig)
	}
}

// handleKreadconfig handles kreadconfig6, kwriteconfig6, kreadconfig5, kwriteconfig5.
// --file, --group, --key, and --value consume the next token and collapse to <val>.
// --type consumes the next token but keeps it literal (structural: bool, int, etc.).
// --delete is a boolean flag kept as-is.
// Remaining positionals collapse to <val>.
func handleKreadconfig(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"--file": true, "--group": true, "--key": true, "--value": true}, Placeholder: "<val>"},
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

		// --type keeps its value literal (structural info)
		if tok == "--type" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: collapse to <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
