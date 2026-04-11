package shellshape

func init() {
	Register("uniq", handleUniq)
}

// handleUniq handles the uniq command.
// Numeric flags (-f/--skip-fields, -s/--skip-chars) consume the next token as N.
// -D/--all-repeated optionally takes a septype keyword (none, prepend, separate).
// Remaining positionals (input_file, output_file) use classifyToken.
func handleUniq(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-f": true, "--skip-fields": true,
		"-s": true, "--skip-chars": true,
		"-w": true, "--check-chars": true,
	}

	septypes := map[string]bool{
		"none": true, "prepend": true, "separate": true,
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

		// -D / --all-repeated: optionally consumes a septype keyword
		if tok == "-D" || tok == "--all-repeated" {
			result = append(result, tok)
			i++
			if i < len(args) && septypes[args[i]] {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
