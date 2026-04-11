package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("du", handleDu)
}

// handleDu handles the du (disk usage) command.
// -d/--max-depth consume a numeric argument (N).
// -B/--block-size consume a numeric argument (N).
// -t/--threshold consumes the next arg as <threshold> (can be e.g. "100M").
// -I/--exclude consume the next arg as <pattern>.
// All other flags are boolean. Positionals are file paths via classifyToken.
func handleDu(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-d": true, "--max-depth": true,
		"-B": true, "--block-size": true,
	}
	thresholdFlags := map[string]bool{
		"-t": true, "--threshold": true,
	}
	patternFlags := map[string]bool{
		"-I": true, "--exclude": true,
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

		if thresholdFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<threshold>")
			continue
		}

		if patternFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<pattern>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: classify as path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
