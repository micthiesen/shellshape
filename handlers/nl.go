package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("nl", handleNl)
}

// handleNl handles the nl (line numbering) command.
// Style flags (-b, -f, -h) have their value preserved for fixed styles (a/t/n)
// but collapse p<expr> to p<pattern>.
// Format flag (-n) has its value preserved (ln/rn/rz).
// String flags (-d, -s) collapse their value to <str>.
// Numeric flags (-i, -l, -v, -w) collapse their value to N.
// -p is boolean. All positionals are file paths.
func handleNl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	styleFlags := map[byte]bool{'b': true, 'f': true, 'h': true}
	formatFlags := map[byte]bool{'n': true}
	stringFlags := map[byte]bool{'d': true, 's': true}
	numericFlags := map[byte]bool{'i': true, 'l': true, 'v': true, 'w': true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Try to match a flag (standalone or fused like -ba, -nrz, -i5)
		if len(tok) >= 2 && tok[0] == '-' && tok[1] != '-' {
			flagChar := tok[1]

			if styleFlags[flagChar] {
				flag := tok[:2]
				val := tok[2:]
				if val == "" {
					result = append(result, flag)
					i++
					if i < len(args) {
						result = append(result, nlStyleValue(args[i]))
						i++
					}
				} else {
					result = append(result, flag, nlStyleValue(val))
					i++
				}
				continue
			}

			if formatFlags[flagChar] {
				flag := tok[:2]
				val := tok[2:]
				if val == "" {
					result = append(result, flag)
					i++
					if i < len(args) {
						result = append(result, args[i])
						i++
					}
				} else {
					result = append(result, flag, val)
					i++
				}
				continue
			}

			if stringFlags[flagChar] {
				flag := tok[:2]
				val := tok[2:]
				result = append(result, flag, "<str>")
				if val == "" {
					i++
					if i < len(args) {
						i++ // skip value
					}
				} else {
					i++
				}
				continue
			}

			if numericFlags[flagChar] {
				flag := tok[:2]
				result = append(result, flag, "N")
				val := tok[2:]
				if val == "" {
					i++
					if i < len(args) {
						i++ // skip value
					}
				} else {
					i++
				}
				continue
			}
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// nlStyleValue normalizes a style value for -b/-f/-h flags.
// Fixed styles (a, t, n) are preserved. Pattern styles (p<expr>) collapse to p<pattern>.
func nlStyleValue(val string) string {
	if strings.HasPrefix(val, "p") && len(val) > 1 {
		return "p<pattern>"
	}
	return val
}
