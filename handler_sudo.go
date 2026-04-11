package shellshape

import "strings"

func init() {
	Register("sudo", handleSudo)
}

// handleSudo handles the sudo command.
// sudo-specific flags are parsed first: -u/-g take a username/group (<val>),
// -p takes a prompt string (<str>), -C takes a file descriptor number (N),
// -r/-t take SELinux role/type (<val>).
// Boolean flags (-A, -b, -E, -e, -H, -h, -i, -K, -k, -l, -n, -P, -S, -s, -V, -v)
// are preserved verbatim.
// Once sudo flags end, remaining tokens are the inner command, classified
// generically via classifyToken.
func handleSudo(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valFlags := map[string]bool{
		"-u": true, "--user": true,
		"-g": true, "--group": true,
		"-r": true, "--role": true,
		"-t": true, "--type": true,
	}
	strFlags := map[string]bool{
		"-p": true, "--prompt": true,
	}
	numericFlags := map[string]bool{
		"-C": true, "--close-from": true,
	}

	var result []string
	parsingFlags := true

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			parsingFlags = false
			i++
			continue
		}

		if parsingFlags {
			if valFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					if isSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "<val>")
					}
					i++
				}
				continue
			}

			if strFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					if isSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "<str>")
					}
					i++
				}
				continue
			}

			if numericFlags[tok] {
				result = append(result, tok)
				i++
				if i < len(args) {
					if isSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "N")
					}
					i++
				}
				continue
			}

			// Long flags with = (e.g. --user=www)
			if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
				eqIdx := strings.Index(tok, "=")
				flagName := tok[:eqIdx]
				if valFlags[flagName] || strFlags[flagName] || numericFlags[flagName] {
					result = append(result, flagName+"=<val>")
					i++
					continue
				}
			}

			if isFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}

			// First non-flag token: stop parsing sudo flags
			parsingFlags = false
		}

		// Inner command tokens: classify generically
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
