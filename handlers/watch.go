package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("watch", handleWatch)
}

// handleWatch handles the watch command.
// POSIX option processing: flags are parsed until the first non-option argument,
// then everything else is the watched command (classified generically via classifyToken).
// -n/--interval and -q/--equexit consume the next token as N.
// -s/--shotsdir consumes the next token as <path>.
// All other flags are boolean.
func handleWatch(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "--interval": true,
		"-q": true, "--equexit": true,
	}
	pathFlags := map[string]bool{
		"-s": true, "--shotsdir": true,
	}

	var result []string
	parsingFlags := true

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			parsingFlags = false
			i++
			continue
		}

		if parsingFlags {
			if numericFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
				continue
			}

			if pathFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
				continue
			}

			if shellshape.IsFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}

			// First non-flag token: stop parsing watch flags
			parsingFlags = false
		}

		// Watched command tokens: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
