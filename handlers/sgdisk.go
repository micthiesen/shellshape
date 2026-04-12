package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("sgdisk", handleSgdisk)
}

func handleSgdisk(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-n": true, "-t": true, "-c": true,
	}
	numericFlags := map[string]bool{
		"-d": true, "-i": true,
	}
	pathFlags := map[string]bool{
		"--backup": true, "--load-backup": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: pathFlags, Placeholder: "<path>"},
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

		// Check flag categories
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Fused long flags (e.g. --backup=/tmp/file)
		if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, norm)
			i++
			continue
		}

		// Boolean flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
