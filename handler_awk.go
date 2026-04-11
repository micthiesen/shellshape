package shellshape

import "strings"

// handleAwk handles awk and gawk.
// -F fs: field separator, next arg becomes <val>. Fused form -F: also accepted.
// -v var=value: variable assignment, next arg becomes <val>.
// -f progfile: program file, next arg becomes <path>. Suppresses positional program.
// First positional (if no -f given) becomes <awk-prog>.
// Remaining positionals use classifyToken.
func handleAwk(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	fsFlags := map[string]bool{"-F": true}
	varFlags := map[string]bool{"-v": true}
	fileFlags := map[string]bool{"-f": true}

	var result []string
	progAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			progAssigned = true
			i++
			continue
		}

		// Fused -F<sep> (e.g. -F: -F, -F\\t)
		if strings.HasPrefix(tok, "-F") && len(tok) > 2 {
			result = append(result, "-F", "<val>")
			i++
			continue
		}

		if fsFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if varFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if fileFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				progAssigned = true
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !progAssigned {
			result = append(result, "<awk-prog>")
			progAssigned = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
