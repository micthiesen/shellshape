package shellshape

import "regexp"

var sectionRE = regexp.MustCompile(`^[1-9][a-z]*$`)

// handleMan handles man, apropos, whatis.
// Flags like -M (manpath), -P (pager), -S (sections), -s, -m, -p consume a value.
// Boolean flags: -a, -d, -f, -h, -k, -K, -o, -t, -w.
// Positionals are kept verbatim: an optional section number (1-9) and page name(s).
func handleMan(_ string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valueFlags := map[string]bool{
		"-M": true, "--manpath": true,
		"-P": true, "--pager": true,
		"-S": true, "--sections": true,
		"-s": true,
		"-m": true,
		"-p": true,
	}

	pathFlags := map[string]bool{
		"-M": true, "--manpath": true,
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

		if valueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if pathFlags[tok] {
					result = append(result, "<path>")
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: keep verbatim (section number or page name)
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
