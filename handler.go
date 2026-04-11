package shellshape

// handlerFunc processes tokens after the executable name.
// subcommand is the subcommand token (empty if none).
type handlerFunc func(subcommand string, tokens []string) []string

var redirectConsumeNext = map[string]bool{
	">": true, ">>": true, "<": true,
	"&>": true, "&>>": true,
	"2>": true, "2>>": true, "1>": true,
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
				redirects = append(redirects, tok, "<path>")
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
