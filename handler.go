package shellshape

import "strings"

// handlerFunc processes tokens after the executable name.
// subcommand is the subcommand token (empty if none).
type handlerFunc func(subcommand string, tokens []string) []string

// flagCategory pairs a set of flag names with the placeholder to use for their values.
type flagCategory struct {
	flags       map[string]bool
	placeholder string
}

// consumeFlagArg appends a flag and its next-token value (collapsed to placeholder)
// to result. Subshell tokens are preserved verbatim.
func consumeFlagArg(tok string, args []string, i int, result []string, placeholder string) ([]string, int) {
	result = append(result, tok)
	i++
	if i < len(args) {
		if isSubshellToken(args[i]) {
			result = append(result, args[i])
		} else {
			result = append(result, placeholder)
		}
		i++
	}
	return result, i
}

// matchFlagCategory checks tok against a list of flag categories and returns
// the matching placeholder. Returns ("", false) if no category matches.
func matchFlagCategory(tok string, categories []flagCategory) (string, bool) {
	for _, cat := range categories {
		if cat.flags[tok] {
			return cat.placeholder, true
		}
	}
	return "", false
}

// consumeFusedFlag handles --flag=value syntax against a list of flag categories.
// Returns the normalized "flag=<placeholder>" string and true if matched,
// or ("", false) if the token doesn't contain '=' or doesn't match any category.
func consumeFusedFlag(tok string, categories []flagCategory) (string, bool) {
	eqIdx := strings.IndexByte(tok, '=')
	if eqIdx < 0 {
		return "", false
	}
	key := tok[:eqIdx]
	for _, cat := range categories {
		if cat.flags[key] {
			return key + "=" + cat.placeholder, true
		}
	}
	return "", false
}

var redirectConsumeNext = map[string]bool{
	">": true, ">>": true, "<": true,
	"&>": true, "&>>": true,
	"2>": true, "2>>": true, "1>": true,
	"<<<": true,
}

var redirectStandalone = map[string]bool{
	"2>&1": true, "1>&2": true,
	"&>/dev/null": true, "2>/dev/null": true,
}

// splitRedirects partitions tokens into non-redirect args and redirect tail.
func splitRedirects(tokens []string) (args, redirects []string) {
	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		if redirectStandalone[tok] {
			redirects = append(redirects, tok)
			i++
			continue
		}

		if redirectConsumeNext[tok] {
			if i+1 < len(tokens) {
				placeholder := "<path>"
				if tok == "<<<" {
					placeholder = "<str>"
				}
				redirects = append(redirects, tok, placeholder)
				i += 2
			} else {
				redirects = append(redirects, tok)
				i++
			}
			continue
		}

		args = append(args, tok)
		i++
	}
	return
}
