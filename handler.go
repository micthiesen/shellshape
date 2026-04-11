package shellshape

import "strings"

// handlerFunc processes tokens after the executable name.
// subcommand is the subcommand token (empty if none).
type HandlerFunc func(subcommand string, tokens []string) []string

// flagCategory pairs a set of flag names with the placeholder to use for their values.
type FlagCategory struct {
	Flags       map[string]bool
	Placeholder string
}

// consumeFlagArg appends a flag and its next-token value (collapsed to placeholder)
// to result. Subshell tokens are preserved verbatim.
func ConsumeFlagArg(tok string, args []string, i int, result []string, placeholder string) ([]string, int) {
	result = append(result, tok)
	i++
	if i < len(args) {
		if IsSubshellToken(args[i]) {
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
func MatchFlagCategory(tok string, categories []FlagCategory) (string, bool) {
	for _, cat := range categories {
		if cat.Flags[tok] {
			return cat.Placeholder, true
		}
	}
	return "", false
}

// consumeFusedFlag handles --flag=value syntax against a list of flag categories.
// Returns the normalized "flag=<placeholder>" string and true if matched,
// or ("", false) if the token doesn't contain '=' or doesn't match any category.
func ConsumeFusedFlag(tok string, categories []FlagCategory) (string, bool) {
	eqIdx := strings.IndexByte(tok, '=')
	if eqIdx < 0 {
		return "", false
	}
	key := tok[:eqIdx]
	for _, cat := range categories {
		if cat.Flags[key] {
			return key + "=" + cat.Placeholder, true
		}
	}
	return "", false
}

var RedirectConsumeNext = map[string]bool{
	">": true, ">>": true, "<": true,
	"&>": true, "&>>": true,
	"2>": true, "2>>": true, "1>": true,
	"<<<": true,
}

var RedirectStandalone = map[string]bool{
	"2>&1": true, "1>&2": true,
	"&>/dev/null": true, "2>/dev/null": true,
}

// splitRedirects partitions tokens into non-redirect args and redirect tail.
func SplitRedirects(tokens []string) (args, redirects []string) {
	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		if RedirectStandalone[tok] {
			redirects = append(redirects, tok)
			i++
			continue
		}

		if RedirectConsumeNext[tok] {
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
