package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("htpasswd", handleHtpasswd)
}

// handleHtpasswd handles the htpasswd command.
// -C (bcrypt cost) consumes the next token as N.
// When -n is present (stdout mode), there is no password file positional.
// Positionals: [passwdfile] username [password]
// passwdfile → <path>, username → <user>, password → <str>.
func handleHtpasswd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Pre-scan for -n flag (may be bundled like -nb, -nbm).
	stdoutMode := false
	for _, tok := range args {
		if tok == "-n" || tok == "--stdout" {
			stdoutMode = true
			break
		}
		// Check for -n bundled in short flags like -nb, -nbm.
		if shellshape.IsFlagToken(tok) && strings.HasPrefix(tok, "-") && !strings.HasPrefix(tok, "--") && strings.ContainsRune(tok, 'n') {
			stdoutMode = true
			break
		}
	}

	numericFlags := map[string]bool{
		"-C": true,
	}

	var result []string
	positional := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			positional++
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals depend on whether -n (stdout mode) is active.
		// With -n:    pos0=<user>, pos1=<str>
		// Without -n: pos0=<path>, pos1=<user>, pos2=<str>
		if stdoutMode {
			switch positional {
			case 0:
				result = append(result, "<user>")
			default:
				result = append(result, "<str>")
			}
		} else {
			switch positional {
			case 0:
				result = append(result, "<path>")
			case 1:
				result = append(result, "<user>")
			default:
				result = append(result, "<str>")
			}
		}
		positional++
		i++
	}

	result = append(result, redirects...)
	return result
}
