package shellshape

import "strings"

func init() {
	Register("make", handleMake, HandlerOptions{HasSubcommands: true})
}

// handleMake handles make arguments after the subcommand (target) extraction.
// Since make is in subcommandExecutables, the first token (typically a target
// name) is already consumed by the framework and kept verbatim.
//
// Path flags (-f, -C, -I, -o, -W and long forms) collapse their arg to <path>.
// Numeric flags (-j, -l and long forms) collapse their arg to N.
// Variable assignments (VAR=value) collapse to VAR=<val>.
// Additional positional targets are kept verbatim.
func handleMake(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-f": true, "--file": true, "--makefile": true,
		"-C": true, "--directory": true,
		"-I": true, "--include-dir": true,
		"-o": true, "--old-file": true, "--assume-old": true,
		"-W": true, "--what-if": true, "--new-file": true, "--assume-new": true,
	}

	numericFlags := map[string]bool{
		"-j": true, "--jobs": true,
		"-l": true, "--load-average": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{numericFlags, "N"},
	}

	var result []string

	// If the subcommand was actually a flag that consumes an argument,
	// the first remaining token is the dangling flag argument.
	i := 0
	if placeholder, ok := matchFlagCategory(subcommand, categories); ok {
		if i < len(args) && !isSubshellToken(args[i]) {
			result = append(result, placeholder)
			i++
		}
	}

	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Long flag with =value
		if norm, ok := consumeFusedFlag(tok, categories); ok {
			result = append(result, norm)
			i++
			continue
		}
		// Unknown fused long flag: collapse value to <val>
		if strings.HasPrefix(tok, "--") {
			if idx := strings.Index(tok, "="); idx >= 0 {
				result = append(result, tok[:idx]+"=<val>")
				i++
				continue
			}
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Variable assignment: VAR=value -> VAR=<val>
		if eqIdx := strings.Index(tok, "="); eqIdx > 0 {
			key := tok[:eqIdx]
			isVar := true
			for _, c := range key {
				if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_') {
					isVar = false
					break
				}
			}
			if isVar {
				result = append(result, key+"=<val>")
				i++
				continue
			}
		}

		// Additional target: keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
