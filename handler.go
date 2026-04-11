package shellshape

// handlerFunc processes tokens after the executable name.
type handlerFunc func(tokens []string) []string

// handler registry - maps executable names to their handler
var handlers = map[string]handlerFunc{
	"echo":   handleEcho,
	"printf": handleEcho,
	"sed":    handleSed,
	"grep":   handleGrep,
	"egrep":  handleGrep,
	"fgrep":  handleGrep,
	"rg":     handleGrep,
	"ag":     handleGrep,
	"ack":    handleGrep,
	"find":    handleFind,
	"sqlite3": handleSqlite3,
	"sqlite":  handleSqlite3,
}

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
