package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("smartctl", handleSmartctl)
}

// handleSmartctl handles smartctl SMART disk monitoring commands.
// Structural flags: -t/--test (test type), -d/--device (device type),
// -l/--log (log type), -n/--nocheck (power mode) — values kept verbatim.
// Boolean flags: -a, -x, -i, -H, -c, -A, -X, --scan, --all, etc.
// The positional is always a device path → <path>.
func handleSmartctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is structural (kept verbatim)
	structuralFlags := map[string]bool{
		"-t": true, "--test": true,
		"-d": true, "--device": true,
		"-l": true, "--log": true,
		"-n": true, "--nocheck": true,
	}

	// Long flags that accept fused =value and should keep value verbatim
	structuralLongFlags := map[string]bool{
		"--test": true, "--device": true, "--log": true, "--nocheck": true,
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle --flag=value for structural flags (keep value verbatim)
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 && strings.HasPrefix(tok, "--") {
			key := tok[:eqIdx]
			if structuralLongFlags[key] {
				result = append(result, tok)
				i++
				continue
			}
		}

		// Structural flags: keep next token verbatim
		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
