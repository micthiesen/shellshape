package shellshape

import "strings"

// handleStow handles GNU stow.
// Path flags (-d, --dir, -t, --target) consume next arg as <path>.
// Value flags (--ignore, --defer, --override) consume next arg as <val>.
// Boolean flags (-D, -R, -S, -n, -v, --adopt, etc.) are preserved.
// Positionals are package names and kept verbatim (they're structural).
// Flags with =value syntax (--dir=X, --ignore=X) collapse the value.
func handleStow(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-d": true, "--dir": true,
		"-t": true, "--target": true,
	}
	valueFlags := map[string]bool{
		"--ignore": true, "--defer": true, "--override": true,
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle --flag=value syntax for path and value flags.
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			eqIdx := strings.Index(tok, "=")
			flagName := tok[:eqIdx]
			if pathFlags[flagName] || valueFlags[flagName] {
				result = append(result, flagName+"=<val>")
				i++
				continue
			}
			// Unknown --flag=value: preserve as-is.
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if valueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: package name, keep verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
