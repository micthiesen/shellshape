package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	shellshape.Register("lscpu", handleLscpu)
}

// handleLscpu handles lscpu (CPU information).
// All flags are boolean. --extended=COLS and --parse=COLS fused forms
// collapse the column list to <val>.
func handleLscpu(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	fusedValueFlags := map[string]bool{
		"-e": true, "--extended": true,
		"-p": true, "--parse": true,
	}

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// Handle --flag=value for fused flags
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 {
			key := tok[:eqIdx]
			if fusedValueFlags[key] {
				result = append(result, key+"=<val>")
				continue
			}
		}

		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}
