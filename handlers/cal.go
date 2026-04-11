package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"cal", "ncal"} {
		shellshape.Register(name, handleCal)
	}
}

// handleCal handles cal and ncal commands.
// Boolean flags (-3, -h, -j, -y, -J, -e, -o, -p, -w, -C, -N) are kept verbatim.
// Flags -A, -B, -m consume the next token as N (numeric).
// Flags -d, -H consume the next token as <date>.
// Flag -s consumes the next token verbatim (country code).
// All positionals (month, year numbers) collapse to N.
func handleCal(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-A": true, "-B": true, "-m": true,
	}
	dateFlags := map[string]bool{
		"-d": true, "-H": true,
	}
	verbatimFlags := map[string]bool{
		"-s": true,
	}
	boolFlags := map[string]bool{
		"-3": true,
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

		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if dateFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<date>")
			continue
		}

		if verbatimFlags[tok] {
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

		// Positional: month or year number
		result = append(result, "N")
		i++
	}

	result = append(result, redirects...)
	return result
}
