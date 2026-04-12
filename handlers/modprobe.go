package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("modprobe", handleModprobe)
}

// handleModprobe handles modprobe arguments.
// Module names are structural (kept verbatim). Module parameters (key=value) → <val>.
// -C → <path>, -d → <path>. Boolean flags: -r, -v, -n, -q, --first-time, etc.
func handleModprobe(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-C": true,
		"-d": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	var result []string
	moduleFound := false
	i := 0

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Path flags
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Other flags: keep verbatim (boolean flags like -r, -v, -n, -q, etc.)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the module name (structural, keep verbatim)
		if !moduleFound {
			result = append(result, tok)
			moduleFound = true
			i++
			continue
		}

		// Subsequent positionals: module parameters (key=value) → <val>
		if strings.Contains(tok, "=") {
			result = append(result, "<val>")
		} else {
			// Non-key=value positional after module: also collapse
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
