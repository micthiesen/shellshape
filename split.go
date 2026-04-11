package shellshape

import "sort"

type segment struct {
	text string
	sep  string
}

var splitOps = []string{"&&", "||", "|&", "|", ";", "\r\n", "\n"}

// splitTopLevel splits a command string on shell operators at the top level only
// (outside quotes, backticks, and subshells). Returns segments with text and separator.
// Newline operators are canonicalized to ";" in the sep field.
func splitTopLevel(command string, operators []string) []segment {
	// Sort operators longest-first
	sorted := make([]string, len(operators))
	copy(sorted, operators)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	var segments []segment
	var cur []byte

	inSingle := false
	inDouble := false
	inBacktick := false
	subshellDepth := 0

	i := 0
	for i < len(command) {
		ch := command[i]

		if inSingle {
			cur = append(cur, ch)
			if ch == '\'' {
				inSingle = false
			}
			i++
			continue
		}

		if inBacktick {
			cur = append(cur, ch)
			if ch == '`' {
				inBacktick = false
			}
			i++
			continue
		}

		if inDouble {
			cur = append(cur, ch)
			if ch == '"' {
				inDouble = false
			}
			i++
			continue
		}

		if ch == '\'' {
			inSingle = true
			cur = append(cur, ch)
			i++
			continue
		}
		if ch == '"' {
			inDouble = true
			cur = append(cur, ch)
			i++
			continue
		}
		if ch == '`' {
			inBacktick = true
			cur = append(cur, ch)
			i++
			continue
		}

		// Track $( for subshell
		if ch == '$' && i+1 < len(command) && command[i+1] == '(' {
			subshellDepth++
			cur = append(cur, ch, command[i+1])
			i += 2
			continue
		}

		if ch == '(' && subshellDepth > 0 {
			subshellDepth++
			cur = append(cur, ch)
			i++
			continue
		}

		if ch == ')' && subshellDepth > 0 {
			subshellDepth--
			cur = append(cur, ch)
			i++
			continue
		}

		if subshellDepth > 0 {
			cur = append(cur, ch)
			i++
			continue
		}

		// Try to match an operator
		matched := false
		for _, op := range sorted {
			if i+len(op) <= len(command) && command[i:i+len(op)] == op {
				sep := op
				if sep == "\n" || sep == "\r\n" {
					sep = ";"
				}
				segments = append(segments, segment{text: string(cur), sep: sep})
				cur = cur[:0]
				i += len(op)
				matched = true
				break
			}
		}
		if matched {
			continue
		}

		cur = append(cur, ch)
		i++
	}

	// Final segment
	segments = append(segments, segment{text: string(cur), sep: ""})

	return segments
}
