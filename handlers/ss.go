package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("ss", handleSs)
}

// handleSs handles the ss (socket statistics) command.
// Value flags (-f, -A, --family, --query, --socket) collapse next arg to <val>.
// Path flags (-D, -F, --diag, --filter) collapse next arg to <path>.
// Keywords "state" and "exclude" collapse the next arg to <state>.
// Keywords "src" and "dst" collapse the next arg to <addr>.
// -4 and -6 are boolean flags (address family selectors).
// Remaining positionals (typically filter expressions) become <filter>.
func handleSs(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

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

	categories := []shellshape.FlagCategory{
		{Flags: valueFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: stateKeywords, Placeholder: "<state>"},
		{Flags: addrKeywords, Placeholder: "<addr>"},
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if digitFlags[tok] || shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
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
		result = append(result, shellshape.ClassifyToken(tok))
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
