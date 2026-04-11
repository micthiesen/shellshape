package shellshape

import (
	"fmt"
	"strings"
)

// shelxSplit splits a string into tokens using POSIX shell-like rules.
// It handles single quotes, double quotes, and escape characters.
// Returns an error on unmatched quotes.
func shelxSplit(s string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	inToken := false
	inSingle := false
	inDouble := false

	i := 0
	for i < len(s) {
		c := s[i]

		if inSingle {
			if c == '\'' {
				inSingle = false
			} else {
				current.WriteByte(c)
			}
			i++
			continue
		}

		if inDouble {
			if c == '"' {
				inDouble = false
			} else if c == '\\' && i+1 < len(s) {
				next := s[i+1]
				// In double quotes, backslash only escapes: $, `, ", \, newline
				if next == '$' || next == '`' || next == '"' || next == '\\' || next == '\n' {
					current.WriteByte(next)
					i += 2
					continue
				}
				current.WriteByte(c)
			} else {
				current.WriteByte(c)
			}
			i++
			continue
		}

		// Outside quotes
		switch c {
		case '\'':
			inSingle = true
			inToken = true
			i++
		case '"':
			inDouble = true
			inToken = true
			i++
		case '\\':
			if i+1 < len(s) {
				current.WriteByte(s[i+1])
				inToken = true
				i += 2
			} else {
				current.WriteByte(c)
				inToken = true
				i++
			}
		case ' ', '\t':
			if inToken {
				tokens = append(tokens, current.String())
				current.Reset()
				inToken = false
			}
			i++
		default:
			current.WriteByte(c)
			inToken = true
			i++
		}
	}

	if inSingle || inDouble {
		return nil, fmt.Errorf("unmatched quote")
	}

	if inToken {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}
