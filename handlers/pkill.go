package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"pkill", "pgrep"} {
		shellshape.Register(name, handlePkill)
	}
}

// handlePkill handles pkill and pgrep commands.
// Signal flags (-9, -HUP, -SIGTERM, --signal NAME) are kept verbatim.
// Boolean match flags (-f, -x, -n, -o, -c, -l, -a, -i, -v) are kept verbatim.
// Value-consuming flags (-u, -g, -G, -P, -s, -t) collapse their argument.
// The positional argument (pattern) becomes <pattern>.
func handlePkill(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	boolFlags := map[string]bool{
		"-f": true, "-x": true, "-n": true, "-o": true,
		"-c": true, "-l": true, "-a": true, "-i": true, "-v": true,
		"--full": true, "--exact": true, "--newest": true, "--oldest": true,
		"--count": true, "--list-name": true, "--list-full": true,
		"--inverse": true, "--lightweight": true,
	}

	userFlags := map[string]bool{
		"-u": true, "--euid": true, "-U": true, "--uid": true,
	}

	valFlags := map[string]bool{
		"-g": true, "-G": true, "-P": true, "-s": true, "-t": true,
		"--pgroup": true, "--group": true, "--parent": true,
		"--session": true, "--terminal": true,
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

		// --signal flag: next token is signal name, keep verbatim
		if tok == "--signal" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i]) // signal name kept verbatim
				i++
			}
			continue
		}

		// Boolean flags
		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// User-consuming flags
		if userFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<user>")
			continue
		}

		// Value-consuming flags
		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Signal as flag: -9, -HUP, -TERM, -SIGKILL, etc.
		if isSignalFlag(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Other flags (unknown, pass through)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: the pattern
		result = append(result, "<pattern>")
		i++
	}

	result = append(result, redirects...)
	return result
}

// isSignalFlag matches -9, -15, -HUP, -SIGTERM, etc.
func isSignalFlag(tok string) bool {
	if !strings.HasPrefix(tok, "-") || len(tok) < 2 {
		return false
	}
	rest := tok[1:]
	// Numeric signal: -9, -15
	allDigits := true
	for _, c := range rest {
		if c < '0' || c > '9' {
			allDigits = false
			break
		}
	}
	if allDigits {
		return true
	}
	// Named signal: -HUP, -TERM, -SIGKILL (all uppercase letters)
	for _, c := range rest {
		if (c < 'A' || c > 'Z') && c != '_' {
			return false
		}
	}
	return true
}
