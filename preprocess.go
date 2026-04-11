package shellshape

import (
	"regexp"
	"strings"
)

// joinContinuations joins backslash-newline line continuations outside of quotes.
// Inside quotes, backslash-newline is left as-is.
func joinContinuations(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	inSingle := false
	inDouble := false

	i := 0
	for i < len(s) {
		ch := s[i]

		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			b.WriteByte(ch)
			i++
		} else if ch == '"' && !inSingle {
			inDouble = !inDouble
			b.WriteByte(ch)
			i++
		} else if ch == '\\' && i+1 < len(s) && s[i+1] == '\n' && !inSingle && !inDouble {
			// Skip the backslash and newline
			i += 2
			// Skip leading whitespace on the next line
			for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
				i++
			}
			// Insert a space if the previous char wasn't whitespace
			if b.Len() > 0 {
				written := b.String()
				last := written[len(written)-1]
				if last != ' ' && last != '\t' && last != '\n' && last != '\r' {
					b.WriteByte(' ')
				}
			}
		} else {
			b.WriteByte(ch)
			i++
		}
	}

	return b.String()
}

var heredocOpenerRe = regexp.MustCompile(`<<-?\s*['"]?([A-Za-z_]\w*)['"]?`)

// collapseHeredocs replaces heredoc bodies with <heredoc>.
// Tracks quote state so << inside strings is ignored.
func collapseHeredocs(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	inSingle := false
	inDouble := false

	i := 0
	for i < len(s) {
		ch := s[i]

		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			b.WriteByte(ch)
			i++
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			b.WriteByte(ch)
			i++
			continue
		}

		if inSingle || inDouble {
			b.WriteByte(ch)
			i++
			continue
		}

		// Check for heredoc opener
		if ch == '<' && i+1 < len(s) && s[i+1] == '<' {
			loc := heredocOpenerRe.FindStringIndex(s[i:])
			if loc != nil && loc[0] == 0 {
				submatch := heredocOpenerRe.FindStringSubmatch(s[i:])
				word := submatch[1]
				openerEnd := i + loc[1]

				// Build closing delimiter regex
				closingRe := regexp.MustCompile(`\n[ \t]*` + regexp.QuoteMeta(word) + `(?:\n|$)`)
				closeLoc := closingRe.FindStringIndex(s[openerEnd:])

				if closeLoc != nil {
					// Find the end of the closing delimiter line
					closeEnd := openerEnd + closeLoc[1]
					// If the match ended at a newline, preserve it
					if closeEnd < len(s) && closeEnd > 0 && s[closeEnd-1] == '\n' {
						closeEnd-- // don't consume trailing newline
					}
					b.WriteString("<heredoc>")
					i = closeEnd
				} else {
					// No closing delimiter found, replace rest of input
					b.WriteString("<heredoc>")
					return b.String()
				}
				continue
			}
		}

		b.WriteByte(ch)
		i++
	}

	return b.String()
}

// stripComments removes bash # comments outside of quotes, backticks, and subshells.
func stripComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	inSingle := false
	inDouble := false
	inBacktick := false
	subshellDepth := 0

	i := 0
	for i < len(s) {
		ch := s[i]

		if inSingle {
			b.WriteByte(ch)
			if ch == '\'' {
				inSingle = false
			}
			i++
			continue
		}

		if inBacktick {
			b.WriteByte(ch)
			if ch == '`' {
				inBacktick = false
			}
			i++
			continue
		}

		if inDouble {
			b.WriteByte(ch)
			if ch == '"' {
				inDouble = false
			}
			i++
			continue
		}

		// Track $( for subshell
		if ch == '$' && i+1 < len(s) && s[i+1] == '(' {
			subshellDepth++
			b.WriteByte(ch)
			i++
			b.WriteByte(s[i])
			i++
			continue
		}

		if ch == '(' && subshellDepth > 0 {
			subshellDepth++
			b.WriteByte(ch)
			i++
			continue
		}

		if ch == ')' && subshellDepth > 0 {
			subshellDepth--
			b.WriteByte(ch)
			i++
			continue
		}

		if subshellDepth > 0 {
			// Inside subshell, treat # as literal
			if ch == '\'' {
				inSingle = true
			} else if ch == '"' {
				inDouble = true
			} else if ch == '`' {
				inBacktick = true
			}
			b.WriteByte(ch)
			i++
			continue
		}

		if ch == '\'' {
			inSingle = true
			b.WriteByte(ch)
			i++
			continue
		}
		if ch == '"' {
			inDouble = true
			b.WriteByte(ch)
			i++
			continue
		}
		if ch == '`' {
			inBacktick = true
			b.WriteByte(ch)
			i++
			continue
		}

		// Check for comment
		if ch == '#' {
			// Comment only if previous char is nil (start) or in set
			isComment := false
			if b.Len() == 0 {
				isComment = true
			} else {
				written := b.String()
				prev := written[len(written)-1]
				if prev == ' ' || prev == '\t' || prev == '\n' || prev == '\r' || prev == ';' || prev == '&' || prev == '|' || prev == '(' {
					isComment = true
				}
			}

			if isComment {
				// Skip to end of line
				for i < len(s) && s[i] != '\n' {
					i++
				}
				// Preserve the newline
				if i < len(s) && s[i] == '\n' {
					b.WriteByte('\n')
					i++
				}
				continue
			}
		}

		b.WriteByte(ch)
		i++
	}

	return b.String()
}
