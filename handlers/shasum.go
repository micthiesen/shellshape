package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"shasum", "md5sum", "sha1sum", "sha224sum", "sha256sum", "sha384sum", "sha512sum"} {
		shellshape.Register(name, handleShasum)
	}
}

// handleShasum handles shasum, md5sum, sha1sum, sha224sum, sha256sum, sha384sum, sha512sum.
// The -a/--algorithm flag consumes the next token and keeps it verbatim (structural).
// All positional arguments are file paths.
func handleShasum(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	algorithmFlags := map[string]bool{"-a": true, "--algorithm": true}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -a / --algorithm: keep the algorithm number verbatim
		if algorithmFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i]) // keep 256, 512, etc. as-is
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a file path for shasum
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
