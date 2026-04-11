package shellshape

func init() {
	Register("jq", handleJq)
}

// handleJq handles jq commands.
// The first positional becomes <filter> (unless -f/--from-file was used).
// --arg/--argjson consume two args: <name> <val>.
// --slurpfile/--rawfile consume two args: <name> <path>.
// -f/--from-file and -L/--library-path consume one arg as <path>.
// --indent consumes one arg as N.
// --args/--jsonargs make all remaining positionals <val>.
// Remaining positionals after the filter use classifyToken.
func handleJq(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	fromFileFlags := map[string]bool{
		"-f": true, "--from-file": true,
	}
	pathFlags := map[string]bool{
		"-L": true, "--library-path": true,
	}
	numericFlags := map[string]bool{
		"--indent": true,
	}
	// Flags that consume two arguments: name + value
	nameValueFlags := map[string]bool{
		"--arg":     true,
		"--argjson": true,
	}
	// Flags that consume two arguments: name + file path
	nameFileFlags := map[string]bool{
		"--slurpfile": true,
		"--rawfile":   true,
	}

	var result []string
	filterAssigned := false
	argsMode := false // after --args or --jsonargs, all positionals are <val>
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
			if nameValueFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					result = append(result, "<name>")
					i++
				}
				if i < len(args) {
					result = append(result, "<val>")
					i++
				}
				continue
			}

			if nameFileFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					result = append(result, "<name>")
					i++
				}
				if i < len(args) {
					result = append(result, "<path>")
					i++
				}
				continue
			}

			if fromFileFlags[tok] {
				result, i = consumeFlagArg(tok, args, i, result, "<path>")
				filterAssigned = true
				continue
			}

			if pathFlags[tok] {
				result, i = consumeFlagArg(tok, args, i, result, "<path>")
				continue
			}

			if numericFlags[tok] {
				result, i = consumeFlagArg(tok, args, i, result, "N")
				continue
			}

			if tok == "--args" || tok == "--jsonargs" {
				result = append(result, tok)
				argsMode = true
				i++
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
		} else if argsMode {
			result = append(result, "<val>")
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
