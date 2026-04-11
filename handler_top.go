package shellshape

import "strings"

func init() {
	Register("top", handleTop)
}

// handleTop handles the top command.
// Flags that consume the next token get specific placeholders:
//
//	-p → <pid>, -u/-U → <user>, -o/-O → <field>
//
// Numeric-consuming flags: -n, -d, -s, -l, -w → N
// Bundled flags like -Hp are detected: if the last letter is a consuming flag,
// the next token is consumed with the appropriate placeholder.
func handleTop(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Single-char flags that consume the next token with a specific placeholder.
	consumingFlags := map[byte]string{
		'p': "<pid>",
		'u': "<user>", 'U': "<user>",
		'o': "<field>", 'O': "<field>",
	}

	// Single-char flags that consume a numeric argument.
	numericFlags := map[byte]bool{
		'n': true, 'd': true, 's': true, 'l': true, 'w': true,
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

		// Handle dash-prefixed flags (single or bundled).
		if strings.HasPrefix(tok, "-") && len(tok) >= 2 && tok[1] != '-' {
			lastChar := tok[len(tok)-1]

			if placeholder, ok := consumingFlags[lastChar]; ok {
				result, i = consumeFlagArg(tok, args, i, result, placeholder)
				continue
			}

			if numericFlags[lastChar] {
				result, i = consumeFlagArg(tok, args, i, result, "N")
				continue
			}
		}

		// Everything else (boolean flags, long flags, etc.) preserved verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
