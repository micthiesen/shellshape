package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("su", handleSu)
}

// handleSu handles the su command.
// -c/--command takes a command string (<code>).
// -s/--shell takes a shell path (<path>).
// -g/--group and -G/--supp-group take a group name (<val>).
// -w/--whitelist-environment takes an env var list (<val>).
// Boolean flags (-, -l, --login, -m, -p, --preserve-environment, -f, --fast)
// are preserved verbatim.
// Positional arguments (username) collapse to <val>.
func handleSu(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	codeFlags := map[string]bool{
		"-c": true, "--command": true,
	}
	pathFlags := map[string]bool{
		"-s": true, "--shell": true,
	}
	valFlags := map[string]bool{
		"-g": true, "--group": true,
		"-G": true, "--supp-group": true,
		"-w": true, "--whitelist-environment": true,
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

		if codeFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<code>")
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Long flags with = (e.g. --command='ls -la')
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			eqIdx := strings.Index(tok, "=")
			flagName := tok[:eqIdx]
			if codeFlags[flagName] || pathFlags[flagName] || valFlags[flagName] {
				result = append(result, flagName+"=<val>")
				i++
				continue
			}
		}

		// Bare "-" is the login shorthand (equivalent to -l)
		if tok == "-" {
			result = append(result, "-")
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: username
		result = shellshape.EmitPositional(result, tok, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
