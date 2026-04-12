package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"wl-copy", "wl-paste"} {
		shellshape.Register(name, handleWlCopy)
	}
}

// handleWlCopy handles wl-copy and wl-paste commands.
// Flags: -o (paste-once), -f (foreground), -c (clear), -p (primary),
// -n (no-newline), -l (list-types). -t/--type takes a MIME type → <val>.
// For wl-paste, -w takes the rest of args as a command (kept verbatim).
// Positional text → <str>.
func handleWlCopy(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-t": true, "--type": true,
	}

	var result []string
	hasStr := false
	i := 0

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			if hasStr {
				result = append(result, "<str>")
				hasStr = false
			}
			result = append(result, tok)
			i++
			continue
		}

		if valFlags[tok] {
			if hasStr {
				result = append(result, "<str>")
				hasStr = false
			}
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// -w (watch mode for wl-paste): rest of args are a command, keep verbatim
		if tok == "-w" || tok == "--watch" {
			if hasStr {
				result = append(result, "<str>")
				hasStr = false
			}
			result = append(result, tok)
			i++
			for i < len(args) {
				result = append(result, args[i])
				i++
			}
			break
		}

		if shellshape.IsFlagToken(tok) {
			if hasStr {
				result = append(result, "<str>")
				hasStr = false
			}
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: text to copy → <str>
		hasStr = true
		i++
	}

	if hasStr {
		result = append(result, "<str>")
	}

	result = append(result, redirects...)
	return result
}
