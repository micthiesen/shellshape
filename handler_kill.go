package shellshape

import "strings"

func init() {
	Register("kill", handleKill)
}

// handleKill handles the kill command.
// Signal flags (-9, -HUP, -SIGTERM, -s NAME) are kept verbatim as structural.
// -l is kept verbatim; if followed by a number it becomes N.
// All positional arguments (PIDs, job specs) collapse to <pid>.
func handleKill(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -- separator: keep it, everything after is positional
		if tok == "--" {
			result = append(result, tok)
			i++
			for i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<pid>")
				}
				i++
			}
			continue
		}

		// -l flag: list signals
		if tok == "-l" {
			result = append(result, "-l")
			i++
			// Optional exit status argument
			if i < len(args) && !isFlagToken(args[i]) {
				result = append(result, "N")
				i++
			}
			continue
		}

		// -s flag: next token is signal name, keep verbatim
		if tok == "-s" {
			result = append(result, "-s")
			i++
			if i < len(args) {
				result = append(result, args[i]) // signal name kept verbatim
				i++
			}
			continue
		}

		// Signal as flag: -9, -HUP, -TERM, -SIGKILL, etc.
		if isFlagToken(tok) || isNumericSignal(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: PID or job spec
		result = append(result, "<pid>")
		i++
	}

	result = append(result, redirects...)
	return result
}

// isNumericSignal matches flags like -9, -15 (signal numbers).
func isNumericSignal(tok string) bool {
	if !strings.HasPrefix(tok, "-") || len(tok) < 2 {
		return false
	}
	for _, c := range tok[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
