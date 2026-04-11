package shellshape

import "strings"

func init() {
	Register("ps", handlePs)
}

// handlePs handles the ps command.
// Flags that consume the next token get specific placeholders:
//
//	-o/-O → <fmt>, -p → <pid>, -u/-U → <user>, -G → <gid>, -g → <grp>, -t → <tty>
//
// Bundled flags like -fu are detected: if the last letter is a consuming flag,
// the next token is consumed with the appropriate placeholder.
// BSD-style option strings (like "aux") are preserved verbatim as structural.
func handlePs(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Maps single-char flags to their placeholder when they consume an argument.
	consumingFlags := map[byte]string{
		'o': "<fmt>", 'O': "<fmt>",
		'p': "<pid>",
		'u': "<user>", 'U': "<user>",
		'G': "<gid>",
		'g': "<grp>",
		't': "<tty>",
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
				continue
			}
		}

		// Exact match for long flags (none for ps currently, but future-proof).
		// Regular flags and BSD-style option strings are preserved verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
