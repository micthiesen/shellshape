package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"xz", "unxz", "xzcat", "lzma", "unlzma"} {
		shellshape.Register(name, handleXz)
	}
}

// handleXz handles xz, unxz, xzcat, lzma, unlzma.
// --threads consumes next token as N. --memlimit consumes next as <val>.
// Compression levels -0 through -9 are preserved verbatim (structural).
// Boolean flags: -d, -z, -k, -f, -v, -t, -c, -l.
// Positionals are file paths.
func handleXz(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"--threads": true, "-T": true}, Placeholder: "N"},
		{Flags: map[string]bool{"--memlimit": true, "--memlimit-compress": true, "--memlimit-decompress": true, "-M": true}, Placeholder: "<val>"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Compression levels: -0 through -9
		if len(tok) == 2 && tok[0] == '-' && tok[1] >= '0' && tok[1] <= '9' {
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
