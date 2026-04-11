package shellshape

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

	categories := []flagCategory{
		{valFlags, "<val>"},
		{strFlags, "<str>"},
		{numericFlags, "N"},
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
			if placeholder, ok := matchFlagCategory(tok, categories); ok {
				result, i = consumeFlagArg(tok, args, i, result, placeholder)
				continue
			}

			// Long flags with = (e.g. --user=www)
			if norm, ok := consumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
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
