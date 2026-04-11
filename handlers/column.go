package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("column", handleColumn)
}

// handleColumn handles the column command.
// Value Flags: -s (delimiter), -c (column width, numeric), -o (output separator),
// -N/--table-columns, -R/--table-right, -H/--table-hide.
// Boolean Flags: -t, -x, -n, -e, -L.
// All positionals are input file paths.
func handleColumn(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-c": true,
	}
	valueFlags := map[string]bool{
		"-s": true, "-o": true,
		"-N": true, "--table-columns": true,
		"-R": true, "--table-right": true,
		"-H": true, "--table-hide": true,
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

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: input file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
