package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("find", handleFind)
}

// handleFind handles the find command.
// Pattern flags (-name, -iname, -path, etc.) make the next arg <pattern>
// unless it's a subshell (preserved verbatim).
// Flags (isFlagToken) are kept verbatim.
// Subshells are preserved verbatim.
// Tokens starting with - but not flags (find primaries) are kept verbatim.
// Remaining positionals use classifyToken.
func handleFind(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	patternFlags := map[string]bool{
		"-name": true, "-iname": true,
		"-path": true, "-ipath": true,
		"-wholename": true, "-iwholename": true,
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

		if patternFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				next := args[i]
				if shellshape.IsSubshellToken(next) {
					result = append(result, next)
				} else {
					result = append(result, "<pattern>")
				}
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Tokens starting with - but not flags (find primaries like -delete, -print, -type, -o, -not)
		if strings.HasPrefix(tok, "-") {
			result = append(result, tok)
			i++
			continue
		}

		// Remaining positionals
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
