package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("dd", handleDd)
}

// handleDd handles the dd command.
// dd uses key=value operands instead of traditional flags.
// Path operands (if=, of=) collapse to key=<path>.
// Numeric operands (bs=, count=, skip=, etc.) collapse to key=N.
// Keyword operands (conv=, iflag=, oflag=, status=, fillchar=) are preserved verbatim.
func handleDd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathKeys := map[string]bool{
		"if": true, "of": true,
	}

	numericKeys := map[string]bool{
		"bs": true, "ibs": true, "obs": true, "cbs": true,
		"count": true, "skip": true, "seek": true,
		"iseek": true, "oseek": true,
		"files": true, "speed": true,
	}

	var result []string

	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if idx := strings.Index(tok, "="); idx > 0 {
			key := tok[:idx]
			val := tok[idx+1:]

			// If the value contains a subshell, preserve it verbatim.
			if strings.Contains(val, "$(") {
				result = append(result, tok)
				continue
			}

			if pathKeys[key] {
				result = append(result, key+"=<path>")
				continue
			}
			if numericKeys[key] {
				result = append(result, key+"=N")
				continue
			}
			// Keyword operands (conv=, iflag=, etc.): preserve verbatim
			result = append(result, tok)
			continue
		}

		// Bare token (unusual for dd, but handle gracefully)
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
