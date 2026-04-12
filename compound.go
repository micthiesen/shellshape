package shellshape

import (
	"strings"
)

// Compound command keywords that start a compound command.
var compoundKeywords = map[string]bool{
	"for": true, "while": true, "until": true, "if": true,
}

// isCompoundCommand checks if the text starts with a compound keyword.
func isCompoundCommand(text string) bool {
	text = strings.TrimSpace(text)
	word := firstWord(text)
	return compoundKeywords[word]
}

// normalizeCompound handles for/while/until/if compound commands.
// The text is the full compound command (e.g., "for f in a b; do echo $f; done").
func normalizeCompound(text string) string {
	text = strings.TrimSpace(text)
	keyword := firstWord(text)

	switch keyword {
	case "for":
		return normalizeFor(text)
	case "while", "until":
		return normalizeWhileUntil(text)
	case "if":
		return normalizeIf(text)
	}
	// Shouldn't reach here if isCompoundCommand was checked first.
	return normalizeSingleCommand(text)
}

// normalizeFor handles:
//
//	for VAR in VALUES; do BODY; done [REDIR]
//	for ((EXPR)); do BODY; done [REDIR]
func normalizeFor(text string) string {
	rest := strings.TrimSpace(text[len("for"):])

	// C-style for loop: for ((EXPR)); do BODY; done
	if strings.HasPrefix(rest, "((") {
		return normalizeCStyleFor(text, rest)
	}

	// Standard for-in loop.
	// Find " in " keyword boundary.
	inIdx := findKeywordInText(rest, "in")
	if inIdx < 0 {
		return normalizeSingleCommand(text)
	}

	varName := strings.TrimSpace(rest[:inIdx])
	afterIn := rest[inIdx+len("in"):]

	// Find "; do" or "\ndo" boundary in the remaining text.
	doIdx := findSemicolonKeyword(afterIn, "do")
	if doIdx < 0 {
		return normalizeSingleCommand(text)
	}

	afterDo := strings.TrimSpace(afterIn[doIdx+len("do"):])

	// Find matching "; done" at depth 0.
	doneIdx := findMatchingClose(afterDo, "done")
	if doneIdx < 0 {
		return normalizeSingleCommand(text)
	}

	bodyText := strings.TrimSpace(afterDo[:doneIdx])
	afterDone := strings.TrimSpace(afterDo[doneIdx+len("done"):])

	// Normalize parts.
	normalizedBody := Normalize(bodyText)

	var parts []string
	parts = append(parts, "for", varName, "in", "<val>+", ";", "do", normalizedBody, ";", "done")
	if afterDone != "" {
		parts = append(parts, normalizeTrailing(afterDone))
	}
	return strings.Join(parts, " ")
}

// normalizeCStyleFor handles for ((EXPR)); do BODY; done.
func normalizeCStyleFor(text, rest string) string {
	// Find the closing )).
	closeIdx := findClosingDoubleParens(rest)
	if closeIdx < 0 {
		return normalizeSingleCommand(text)
	}

	afterParens := rest[closeIdx+2:]

	// Find "; do" after the parens.
	doIdx := findSemicolonKeyword(afterParens, "do")
	if doIdx < 0 {
		return normalizeSingleCommand(text)
	}
	afterDo := strings.TrimSpace(afterParens[doIdx+len("do"):])

	// Find "; done".
	doneIdx := findMatchingClose(afterDo, "done")
	if doneIdx < 0 {
		return normalizeSingleCommand(text)
	}

	bodyText := strings.TrimSpace(afterDo[:doneIdx])
	afterDone := strings.TrimSpace(afterDo[doneIdx+len("done"):])

	normalizedBody := Normalize(bodyText)

	var parts []string
	parts = append(parts, "for", "((<expr>))", ";", "do", normalizedBody, ";", "done")
	if afterDone != "" {
		parts = append(parts, normalizeTrailing(afterDone))
	}
	return strings.Join(parts, " ")
}

// normalizeWhileUntil handles:
//
//	while COND; do BODY; done [REDIR]
//	until COND; do BODY; done [REDIR]
func normalizeWhileUntil(text string) string {
	keyword := firstWord(text)
	rest := strings.TrimSpace(text[len(keyword):])

	// Find "; do".
	doIdx := findSemicolonKeyword(rest, "do")
	if doIdx < 0 {
		return normalizeSingleCommand(text)
	}

	condText := strings.TrimSpace(rest[:doIdx])
	afterDo := strings.TrimSpace(rest[doIdx+len("do"):])

	// Find "; done".
	doneIdx := findMatchingClose(afterDo, "done")
	if doneIdx < 0 {
		return normalizeSingleCommand(text)
	}

	bodyText := strings.TrimSpace(afterDo[:doneIdx])
	afterDone := strings.TrimSpace(afterDo[doneIdx+len("done"):])

	normalizedCond := Normalize(condText)
	normalizedBody := Normalize(bodyText)

	var parts []string
	parts = append(parts, keyword, normalizedCond, ";", "do", normalizedBody, ";", "done")
	if afterDone != "" {
		parts = append(parts, normalizeTrailing(afterDone))
	}
	return strings.Join(parts, " ")
}

// normalizeIf handles:
//
//	if COND; then BODY; [elif COND; then BODY;] [else BODY;] fi [REDIR]
func normalizeIf(text string) string {
	rest := strings.TrimSpace(text[len("if"):])

	// Find "; then".
	thenIdx := findSemicolonKeyword(rest, "then")
	if thenIdx < 0 {
		return normalizeSingleCommand(text)
	}

	condText := strings.TrimSpace(rest[:thenIdx])
	afterThen := strings.TrimSpace(rest[thenIdx+len("then"):])

	normalizedCond := Normalize(condText)

	var parts []string
	parts = append(parts, "if", normalizedCond, ";", "then")

	// Now scan for elif/else/fi at depth 0.
	remaining := afterThen
	for {
		// Try to find the next keyword boundary: elif, else, or fi.
		nextKw, nextIdx := findNextIfKeyword(remaining)
		if nextKw == "" || nextIdx < 0 {
			// Malformed — fall back.
			return normalizeSingleCommand(text)
		}

		bodyText := strings.TrimSpace(remaining[:nextIdx])
		normalizedBody := Normalize(bodyText)
		parts = append(parts, normalizedBody, ";")

		switch nextKw {
		case "fi":
			afterFi := strings.TrimSpace(remaining[nextIdx+len("fi"):])
			parts = append(parts, "fi")
			if afterFi != "" {
				parts = append(parts, normalizeTrailing(afterFi))
			}
			return strings.Join(parts, " ")

		case "elif":
			remaining = strings.TrimSpace(remaining[nextIdx+len("elif"):])
			// Find "; then" after elif.
			thenIdx2 := findSemicolonKeyword(remaining, "then")
			if thenIdx2 < 0 {
				return normalizeSingleCommand(text)
			}
			elifCond := strings.TrimSpace(remaining[:thenIdx2])
			remaining = strings.TrimSpace(remaining[thenIdx2+len("then"):])
			parts = append(parts, "elif", Normalize(elifCond), ";", "then")

		case "else":
			remaining = strings.TrimSpace(remaining[nextIdx+len("else"):])
			parts = append(parts, "else")
		}
	}
}

// normalizeTrailing normalizes redirect tokens that appear after done/fi/esac.
func normalizeTrailing(text string) string {
	tokens, err := shelxSplit(text)
	if err != nil || len(tokens) == 0 {
		return text
	}
	var result []string
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if RedirectConsumeNext[tok] && i+1 < len(tokens) {
			placeholder := "<path>"
			if tok == "<<<" {
				placeholder = "<str>"
			}
			result = append(result, tok, placeholder)
			i++
		} else if RedirectStandalone[tok] {
			result = append(result, tok)
		} else {
			result = append(result, ClassifyToken(tok))
		}
	}
	return strings.Join(result, " ")
}

// --- Scanning helpers ---

// firstWord returns the first whitespace-delimited word in text.
func firstWord(text string) string {
	fields := strings.Fields(text)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

// findKeywordInText finds a keyword at a word boundary in text,
// returning the index of the keyword start, or -1.
func findKeywordInText(text, keyword string) int {
	for i := 0; i <= len(text)-len(keyword); i++ {
		if matchKeywordAt(text, i, keyword) {
			return i
		}
	}
	return -1
}

// findSemicolonKeyword finds "; KEYWORD" or ";\nKEYWORD" at compound-depth 0,
// returning the index of the keyword start (after the semicolon and whitespace).
// It respects quotes and paren depth to avoid matching inside nested constructs.
func findSemicolonKeyword(text, keyword string) int {
	inSingle := false
	inDouble := false
	parenDepth := 0
	compoundD := 0

	for i := 0; i < len(text); i++ {
		ch := text[i]

		if inSingle {
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		if ch == '\'' {
			inSingle = true
			continue
		}
		if ch == '"' {
			inDouble = true
			continue
		}
		if ch == '(' {
			parenDepth++
			continue
		}
		if ch == ')' && parenDepth > 0 {
			parenDepth--
			continue
		}
		if parenDepth > 0 {
			continue
		}

		// Track nested compound depth.
		for _, kw := range compoundOpeners {
			if matchKeywordAt(text, i, kw) {
				compoundD++
				break
			}
		}
		for _, kw := range compoundClosers {
			if matchKeywordAt(text, i, kw) {
				compoundD--
				break
			}
		}

		if compoundD > 0 {
			continue
		}

		// Look for "; keyword" at this position.
		if ch == ';' {
			// Skip whitespace after semicolon.
			j := i + 1
			for j < len(text) && (text[j] == ' ' || text[j] == '\t' || text[j] == '\n') {
				j++
			}
			if matchKeywordAt(text, j, keyword) {
				return j
			}
		}
	}
	return -1
}

// findMatchingClose finds the closing keyword (done/fi/esac) at depth 0,
// preceded by "; " or at the end after "; ". Returns the index of the keyword start.
func findMatchingClose(text, closeKW string) int {
	inSingle := false
	inDouble := false
	parenDepth := 0
	depth := 0 // nested compound depth within the body

	for i := 0; i < len(text); i++ {
		ch := text[i]

		if inSingle {
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		if ch == '\'' {
			inSingle = true
			continue
		}
		if ch == '"' {
			inDouble = true
			continue
		}
		if ch == '(' {
			parenDepth++
			continue
		}
		if ch == ')' && parenDepth > 0 {
			parenDepth--
			continue
		}
		if parenDepth > 0 {
			continue
		}

		// Track nested compounds within the body.
		for _, kw := range compoundOpeners {
			if matchKeywordAt(text, i, kw) {
				depth++
				break
			}
		}

		// Check for closing keyword.
		if matchKeywordAt(text, i, closeKW) {
			if depth == 0 {
				// Find the "; " before the keyword and return the keyword position.
				return i
			}
			depth--
			continue
		}
	}
	return -1
}

// findNextIfKeyword finds the next elif/else/fi keyword at depth 0
// in the remaining text. Returns the keyword name and its start index.
func findNextIfKeyword(text string) (string, int) {
	inSingle := false
	inDouble := false
	parenDepth := 0
	depth := 0

	for i := 0; i < len(text); i++ {
		ch := text[i]

		if inSingle {
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		if ch == '\'' {
			inSingle = true
			continue
		}
		if ch == '"' {
			inDouble = true
			continue
		}
		if ch == '(' {
			parenDepth++
			continue
		}
		if ch == ')' && parenDepth > 0 {
			parenDepth--
			continue
		}
		if parenDepth > 0 {
			continue
		}

		// Track nested if depth.
		if matchKeywordAt(text, i, "if") {
			depth++
			continue
		}
		if matchKeywordAt(text, i, "fi") {
			if depth == 0 {
				return "fi", i
			}
			depth--
			continue
		}

		// elif and else only matter at depth 0.
		if depth == 0 {
			// Check for "; keyword" pattern.
			if ch == ';' {
				j := i + 1
				for j < len(text) && (text[j] == ' ' || text[j] == '\t' || text[j] == '\n') {
					j++
				}
				for _, kw := range []string{"elif", "else"} {
					if matchKeywordAt(text, j, kw) {
						return kw, j
					}
				}
			}
		}
	}
	return "", -1
}

// findClosingDoubleParens finds the )) that closes (( in rest.
// rest starts with "((". Returns the index of the first ) of )).
func findClosingDoubleParens(rest string) int {
	if !strings.HasPrefix(rest, "((") {
		return -1
	}
	depth := 0
	inSingle := false
	inDouble := false
	for i := 0; i < len(rest); i++ {
		ch := rest[i]
		if inSingle {
			if ch == '\'' {
				inSingle = false
			}
			continue
		}
		if inDouble {
			if ch == '"' {
				inDouble = false
			}
			continue
		}
		if ch == '\'' {
			inSingle = true
			continue
		}
		if ch == '"' {
			inDouble = true
			continue
		}
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
			if depth == 0 {
				return i - 1 // index of first ) in ))
			}
		}
	}
	return -1
}
