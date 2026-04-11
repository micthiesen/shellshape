package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("nslookup", handleNslookup)
}

// handleNslookup handles the nslookup DNS lookup utility.
// nslookup uses -flag=value syntax for options with arguments.
// Type flags (-type=, -query=, -querytype=) preserve the type value verbatim.
// Numeric flags (-timeout=, -retry=, -port=) collapse the value to N.
// Boolean flags (-debug, -vc, etc.) are preserved as-is.
// First positional becomes <name>, second becomes <server>.
func handleNslookup(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	typeFlags := map[string]bool{
		"-type": true, "-query": true, "-querytype": true,
	}
	numericFlags := map[string]bool{
		"-timeout": true, "-retry": true, "-port": true,
	}

	var result []string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		// Handle -key=value flags.
		if strings.HasPrefix(tok, "-") && strings.Contains(tok, "=") {
			eqIdx := strings.Index(tok, "=")
			key := tok[:eqIdx]
			val := tok[eqIdx+1:]

			if typeFlags[key] {
				// Preserve type value verbatim.
				result = append(result, key+"="+val)
			} else if numericFlags[key] {
				result = append(result, key+"=N")
			} else {
				// Unknown -key=value, preserve verbatim.
				result = append(result, tok)
			}
			i++
			continue
		}

		// Boolean flags.
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first is <name>, second is <server>.
		positionalCount++
		if positionalCount == 1 {
			result = append(result, "<name>")
		} else {
			result = append(result, "<server>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
