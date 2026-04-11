package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("passwd", handlePasswd)
}

// handlePasswd handles the passwd command.
// The username positional collapses to <user>.
// Numeric flags (-x, -n, -w, -i and long forms) make the next arg N.
// Path flags (-R, -P and long forms) make the next arg <path>.
// Repository flag (-r/--repository) keeps its argument verbatim.
// All other flags are boolean and kept as-is.
func handlePasswd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-x": true, "--maxdays": true,
		"-n": true, "--mindays": true,
		"-w": true, "--warndays": true,
		"-i": true, "--inactive": true,
	}

	pathFlags := map[string]bool{
		"-R": true, "--root": true,
		"-P": true, "--prefix": true,
	}

	repoFlags := map[string]bool{
		"-r": true, "--repository": true,
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

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if repoFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: username
		result = append(result, "<user>")
		i++
	}

	result = append(result, redirects...)
	return result
}
