package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"fprintd-enroll", "fprintd-verify", "fprintd-list", "fprintd-delete"} {
		shellshape.Register(name, handleFprintd)
	}
}

// handleFprintd handles fprintd-enroll, fprintd-verify, fprintd-list, fprintd-delete.
// -f <finger> is structural (finger name like right-index-finger is kept verbatim).
// Positional username collapses to <val>.
func handleFprintd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// -f is structural: finger name kept verbatim
	structuralFlags := map[string]bool{
		"-f": true,
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

		if structuralFlags[tok] {
			// Keep next token verbatim (finger name)
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
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
