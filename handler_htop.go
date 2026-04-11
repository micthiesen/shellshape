package shellshape

func init() {
	Register("htop", handleHtop)
}

// handleHtop handles the htop command.
// Numeric flags (-d, -H) collapse their argument to N.
// Value flags (-p → <pid>, -s → <col>, -u → <user>, -F → <filter>) collapse
// their argument to a descriptive placeholder.
// All other flags are boolean and preserved verbatim.
func handleHtop(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-d": true, "--delay": true,
		"-H": true, "--highlight-changes": true,
	}

	valueFlags := map[string]string{
		"-p": "<pid>", "--pid": "<pid>",
		"-s": "<col>", "--sort-key": "<col>",
		"-u": "<user>", "--user": "<user>",
		"-F": "<filter>", "--filter": "<filter>",
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

		if placeholder, ok := valueFlags[tok]; ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Boolean flags and anything else preserved verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
