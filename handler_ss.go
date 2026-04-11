package shellshape

import "strings"

func init() {
	Register("ss", handleSs)
}

// handleSs handles the ss (socket statistics) command.
// Value flags (-f, -A, --family, --query, --socket) collapse next arg to <val>.
// Path flags (-D, -F, --diag, --filter) collapse next arg to <path>.
// Keywords "state" and "exclude" collapse the next arg to <state>.
// Keywords "src" and "dst" collapse the next arg to <addr>.
// -4 and -6 are boolean flags (address family selectors).
// Remaining positionals (typically filter expressions) become <filter>.
func handleSs(_ string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valueFlags := map[string]bool{
		"-f": true, "--family": true,
		"-A": true, "--query": true, "--socket": true,
	}

	pathFlags := map[string]bool{
		"-D": true, "--diag": true,
		"-F": true, "--filter": true,
	}

	// Keywords that consume the next token.
	stateKeywords := map[string]bool{
		"state": true, "exclude": true,
	}
	addrKeywords := map[string]bool{
		"src": true, "dst": true,
	}

	digitFlags := map[string]bool{
		"-4": true, "-6": true,
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

		if valueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
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

		if digitFlags[tok] || isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// State/exclude keyword: next token is a state name.
		if stateKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<state>")
				}
				i++
			}
			continue
		}

		// Address keyword: next token is an address.
		if addrKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<addr>")
				}
				i++
			}
			continue
		}

		// Any remaining positional is a filter expression.
		// This covers quoted expressions like 'sport = :80' (which arrive
		// as a single token with spaces) and bare filter tokens.
		if isFilterExpression(tok) {
			result = append(result, "<filter>")
			i++
			continue
		}

		// Fallback: classify as generic.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// isFilterExpression returns true if the token looks like an ss filter
// expression. Filter expressions typically contain spaces (from quoted
// strings) or reference socket fields like sport/dport.
func isFilterExpression(tok string) bool {
	// Tokens with spaces are quoted filter expressions.
	if strings.Contains(tok, " ") {
		return true
	}
	return false
}
