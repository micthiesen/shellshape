package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("stow", handleStow)
}

// handleStow handles GNU stow.
// Path flags (-d, --dir, -t, --target) consume next arg as <path>.
// Value flags (--ignore, --defer, --override) consume next arg as <val>.
// Boolean flags (-D, -R, -S, -n, -v, --adopt, etc.) are preserved.
// Positionals are package names, collapsed to <pkg>.
// Flags with =value syntax (--dir=X, --ignore=X) collapse the value.
func handleStow(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

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

		if shellshape.IsSubshellToken(tok) {
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
			// Unknown --flag=value: normalize numeric values.
			val := tok[eqIdx+1:]
			if shellshape.NumberRE.MatchString(val) {
				result = append(result, flagName+"=N")
			} else {
				result = append(result, tok)
			}
			i++
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: package name, collapse to <pkg>.
		result = append(result, "<pkg>")
		i++
	}

	result = append(result, redirects...)
	return result
}
