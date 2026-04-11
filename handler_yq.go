package shellshape

// handleYq handles yq (mikefarah/yq) commands.
// The first positional becomes <filter> (the yq expression).
// Subcommands eval/eval-all/ea are passed through verbatim.
// -o/--output-format, -p/--input-format consume one arg as <val>.
// -I/--indent consumes one arg as N.
// Remaining positionals after the filter use classifyToken.
func handleYq(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valueFlags := map[string]bool{
		"-o": true, "--output-format": true,
		"-p": true, "--input-format": true,
		"--xml-attribute-prefix": true,
		"--xml-content-name":     true,
	}
	numericFlags := map[string]bool{
		"-I": true, "--indent": true,
	}
	subcommands := map[string]bool{
		"eval": true, "eval-all": true, "ea": true,
	}

	var result []string
	filterAssigned := false
	seenDoubleDash := false

	i := 0
	for i < len(args) {
		tok := args[i]

		// Handle -- separator
		if tok == "--" && !seenDoubleDash {
			result = append(result, tok)
			seenDoubleDash = true
			i++
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !filterAssigned {
				filterAssigned = true
			}
			i++
			continue
		}

		if !seenDoubleDash {
			// Pass subcommands through verbatim (only before filter)
			if !filterAssigned && subcommands[tok] {
				result = append(result, tok)
				i++
				continue
			}

			if valueFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					result = append(result, "<val>")
					i++
				}
				continue
			}

			if numericFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					result = append(result, "N")
					i++
				}
				continue
			}

			if isFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}
		}

		// Positional argument
		if !filterAssigned {
			result = append(result, "<filter>")
			filterAssigned = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
