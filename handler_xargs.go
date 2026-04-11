package shellshape

import "strings"

func init() {
	Register("xargs", handleXargs)
}

// handleXargs handles the xargs command.
// xargs flags (-I, -J, -E take strings; -n, -P, -L, -s, -R, -S take numbers)
// are normalized with appropriate placeholders. Boolean flags (-0, -o, -p, -r,
// -t, -x) are preserved verbatim.
// The first non-flag positional after xargs flags is the utility name (preserved).
// Remaining tokens after the utility are treated as utility arguments: flags are
// preserved, subshells are kept verbatim, and other positionals collapse to <arg>.
func handleXargs(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	stringFlags := map[string]bool{
		"-I": true, "-J": true, "-E": true,
	}
	numericFlags := map[string]bool{
		"-n": true, "-P": true, "-L": true,
		"-s": true, "-R": true, "-S": true,
	}

	var result []string
	utilityFound := false

	i := 0
	for i < len(args) {
		tok := args[i]

		// Once we've found the utility, everything else is utility args.
		if utilityFound {
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else if isFlagToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		// Still in xargs flags territory.
		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if stringFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if numericFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// -0 is a boolean flag but isFlagToken doesn't catch it (digit after dash).
		if tok == "-0" {
			result = append(result, tok)
			i++
			continue
		}

		// Long flags with = (e.g. --max-args=5)
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			result = append(result, tok[:strings.Index(tok, "=")]+"=<val>")
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First non-flag positional is the utility name.
		result = append(result, tok)
		utilityFound = true
		i++
	}

	result = append(result, redirects...)
	return result
}
