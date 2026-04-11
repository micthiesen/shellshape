package shellshape

func init() {
	Register("vitest", handleVitest)
}

// handleVitest handles vitest test runner commands.
// Pattern flags (-t, --testNamePattern, --grep) make the next arg <pattern>.
// Path flags (--config, -c, --root, --dir, --outputFile) make the next arg <path>.
// Value flags (--reporter, --environment, --pool) make the next arg <val>.
// Numeric flags (--bail, --retry, --maxConcurrency, --minWorkers, --maxWorkers) make the next arg N.
// Remaining positionals use classifyToken.
func handleVitest(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	patternFlags := map[string]bool{
		"-t": true, "--testNamePattern": true, "--grep": true,
	}
	pathFlags := map[string]bool{
		"--config": true, "-c": true, "--root": true,
		"--dir": true, "--outputFile": true,
	}
	valueFlags := map[string]bool{
		"--reporter": true, "--environment": true, "--pool": true,
	}
	numericFlags := map[string]bool{
		"--bail": true, "--retry": true, "--maxConcurrency": true,
		"--minWorkers": true, "--maxWorkers": true,
	}

	subcommands := map[string]bool{
		"run": true, "bench": true, "watch": true,
	}

	var result []string

	i := 0

	// Detect leading subcommand.
	if i < len(args) && subcommands[args[i]] {
		result = append(result, args[i])
		i++
	}
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if patternFlags[tok] {
			result = append(result, tok)
			i++
			// Consume all non-flag tokens as the pattern (quoted strings get split by shlex)
			for i < len(args) && !isFlagToken(args[i]) && !isSubshellToken(args[i]) {
				i++
			}
			result = append(result, "<pattern>")
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

		// Positional
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
