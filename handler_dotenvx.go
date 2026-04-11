package shellshape

// handleDotenvx handles dotenvx and dotenv (dotenv-cli).
// Both are env var loaders that wrap another command.
// Flags before "--" belong to dotenvx/dotenv; after "--" everything is the
// wrapped command and gets generic classification.
// Path flags: -f/--env-file, -fv/--env-vault-file, -e (dotenv-cli).
// Value flags: -e/--env (dotenvx), --convention.
// Boolean flags: -o/--overload/--override, --no-expand, -p.
func handleDotenvx(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a file path.
	pathFlags := map[string]bool{
		"-f": true, "--env-file": true,
		"-fv": true, "--env-vault-file": true,
	}

	// Flags whose next token is an opaque value.
	valueFlags := map[string]bool{
		"-e": true, "--env": true,
		"--convention": true,
	}

	var result []string
	i := 0

	// Process dotenvx/dotenv flags until "--" or end of args.
	for i < len(args) {
		tok := args[i]

		// "--" separates dotenvx flags from the wrapped command.
		if tok == "--" {
			result = append(result, "--")
			i++
			// Everything after "--" is the wrapped command.
			for i < len(args) {
				result = append(result, classifyToken(args[i]))
				i++
			}
			break
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		if valueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional before "--": classify generically.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
