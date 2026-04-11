package shellshape

func init() {
	Register("column", handleColumn)
}

// handleColumn handles the column command.
// Value flags: -s (delimiter), -c (column width, numeric), -o (output separator),
// -N/--table-columns, -R/--table-right, -H/--table-hide.
// Boolean flags: -t, -x, -n, -e, -L.
// All positionals are input file paths.
func handleColumn(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

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

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: input file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
