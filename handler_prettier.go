package shellshape

func init() {
	Register("prettier", handlePrettier)
}

// handlePrettier handles the prettier code formatter.
// Flags like --config, --ignore-path take a path argument.
// Flags like --parser, --trailing-comma, --arrow-parens take a value argument.
// Flags like --tab-width, --print-width, --range-start, --range-end take a numeric argument.
// Boolean flags (--write, --check, --single-quote, etc.) are kept verbatim.
// All positionals are file paths or glob patterns, classified via classifyToken.
func handlePrettier(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"--config":         true,
		"--ignore-path":    true,
		"--cache-location": true,
	}

	valueFlags := map[string]bool{
		"--parser":         true,
		"--trailing-comma": true,
		"--arrow-parens":   true,
		"--end-of-line":    true,
		"--log-level":      true,
		"--plugin":         true,
	}

	numericFlags := map[string]bool{
		"--tab-width":   true,
		"--print-width": true,
		"--range-start": true,
		"--range-end":   true,
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

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
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

		// Positional: file path or glob pattern
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
