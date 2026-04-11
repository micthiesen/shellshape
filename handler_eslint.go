package shellshape

func init() {
	Register("eslint", handleEslint)
}

// handleEslint handles the eslint linter.
// Path flags (--config/-c, --ignore-path, --output-file/-o, --cache-location)
// collapse their argument to <path>.
// Value flags (--ext, --rule, --format/-f) collapse their argument to <val>.
// Numeric flags (--max-warnings) collapse their argument to N.
// Boolean flags (--fix, --cache, --quiet, etc.) are preserved verbatim.
// Positionals are file/directory paths classified via classifyToken.
func handleEslint(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"--config": true, "-c": true,
		"--ignore-path": true,
		"--output-file": true, "-o": true,
		"--cache-location": true,
	}
	valFlags := map[string]bool{
		"--ext":    true,
		"--rule":   true,
		"--format": true, "-f": true,
	}
	numericFlags := map[string]bool{
		"--max-warnings": true,
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

		if valFlags[tok] {
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

		// Positional: file/directory path or glob pattern
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
