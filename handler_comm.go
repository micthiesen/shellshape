package shellshape

import "regexp"

func init() {
	Register("comm", handleComm)
}

var commFlagRE = regexp.MustCompile(`^-[123]+i?$`)

// handleComm handles the comm command.
// Boolean flags: -1, -2, -3 (suppress columns), -i (case insensitive).
// Flags -1/-2/-3 are digit flags that isFlagToken won't catch, so they
// are matched explicitly via commFlagRE (also handles fused forms like -12, -23i).
// Exactly two positional args, both file paths.
func handleComm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	for i := 0; i < len(args); i++ {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// Digit flags like -1, -2, -3, -12, -123, -23i
		if commFlagRE.MatchString(tok) {
			result = append(result, tok)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// Positional: stdin dash or file path
		if tok == "-" {
			result = append(result, "-")
		} else {
			result = append(result, classifyToken(tok))
		}
	}

	result = append(result, redirects...)
	return result
}
