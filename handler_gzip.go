package shellshape

import "strings"

// handleGzip handles gzip, gunzip, zcat.
// -S/--suffix consumes the next token as <suffix>.
// Compression level flags -1 through -9 are preserved verbatim.
// All positionals are file paths collapsed via classifyToken.
func handleGzip(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	suffixFlags := map[string]bool{"-S": true, "--suffix": true}

	var result []string
	pathCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if suffixFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<suffix>")
				i++
			}
			continue
		}

		// Compression level flags: -1 through -9
		if len(tok) == 2 && tok[0] == '-' && tok[1] >= '1' && tok[1] <= '9' {
			result = append(result, tok)
			i++
			continue
		}

		// --suffix=value
		if strings.HasPrefix(tok, "--suffix=") {
			result = append(result, "--suffix=<suffix>")
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		pathCount++
		if pathCount == 1 {
			result = append(result, classifyToken(tok))
		} else if pathCount == 2 {
			// Replace the first <path> with <path>+
			for j := len(result) - 1; j >= 0; j-- {
				if result[j] == "<path>" {
					result[j] = "<path>+"
					break
				}
			}
		}
		// 3+ paths: already collapsed
		i++
	}

	result = append(result, redirects...)
	return result
}
