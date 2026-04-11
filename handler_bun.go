package shellshape

// handleBun handles the bun command.
// The subcommand (run, test, install, build, x) is already extracted by
// normalizeSingleCommand. This handler processes the remaining tokens after
// the subcommand.
//
// Path flags (--env-file, --config, --outdir) consume the next token as <path>.
// Value flags (--target, --format) consume the next token as <val>.
// Numeric flags (--timeout, --bail, --port) consume the next token as N.
// Boolean flags (--watch, --hot, --frozen-lockfile, --production) are kept verbatim.
// Remaining positionals use classifyToken.
func handleBun(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"--env-file": true, "--config": true, "--outdir": true,
		"--outfile": true,
	}

	valueFlags := map[string]bool{
		"--target": true, "--format": true,
	}

	numericFlags := map[string]bool{
		"--timeout": true, "--bail": true, "--port": true,
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
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
