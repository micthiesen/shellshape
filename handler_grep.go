package shellshape

import "regexp"

var fusedNumericRE = regexp.MustCompile(`^-([ABCm])(\d+)$`)

// handleGrep handles grep, egrep, fgrep, rg, ag, ack.
// Pattern flags (-e, --regexp) make the next arg <pattern>.
// File flags (-f, --file) make the next arg <path>.
// Numeric flags (-A, -B, -C, -m, and long forms) make the next arg N.
// Fused numeric flags like -A3 are split into -A and N.
// First positional (if no pattern assigned) becomes <pattern>.
// Remaining positionals use classifyToken.
func handleGrep(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	patternFlags := map[string]bool{"-e": true, "--regexp": true}
	fileFlags := map[string]bool{"-f": true, "--file": true}
	numericFlags := map[string]bool{
		"-A": true, "-B": true, "-C": true, "-m": true,
		"--max-count": true, "--after-context": true,
		"--before-context": true, "--context": true,
	}

	var result []string
	patternAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			patternAssigned = true
			i++
			continue
		}

		if patternFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<pattern>")
				patternAssigned = true
				i++
			}
			continue
		}

		if fileFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		// Fused numeric: -A3, -B10, etc.
		if m := fusedNumericRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "N")
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !patternAssigned {
			result = append(result, "<pattern>")
			patternAssigned = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
