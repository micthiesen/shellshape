package shellshape

import "regexp"

var portRE = regexp.MustCompile(`^\d+(-\d+)?$`)

// handleNc handles nc, netcat, and ncat.
// Numeric flags (-p, -w, -i, -G, -H, -I, -J, -L) collapse next arg to N.
// Value flags (-s, -X, -x, -b) collapse next arg to <val>.
// Path flags (-e) collapse next arg to <path>.
// Boolean flags (-l, -u, -v, -z, -k, -n, -4, -6, etc.) are preserved.
// Positionals: numeric/port-range tokens become N, others become <host>.
func handleNc(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-p": true, "-w": true, "-i": true,
		"-G": true, "-H": true, "-I": true,
		"-J": true, "-L": true,
	}

	valueFlags := map[string]bool{
		"-s": true, "-X": true, "-x": true, "-b": true,
	}

	pathFlags := map[string]bool{
		"-e": true,
	}

	// Flags that start with a digit and aren't caught by isFlagToken
	digitFlags := map[string]bool{
		"-4": true, "-6": true,
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

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
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

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if digitFlags[tok] || isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: numeric/port-range tokens are ports (N), others are hosts
		if portRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, "<host>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
