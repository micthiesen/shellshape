package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"less", "more"} {
		shellshape.Register(name, handleLess)
	}
}

// handleLess handles less and more.
// Most flags are boolean. A few flags consume the next token:
//   - Numeric: -b, -h, -j, -x, -y, -z (buffer, scroll, target, tabs, scroll limit, window)
//   - Path: -k, -o, -O (lesskey file, log file)
//   - Value: -p, -t (search pattern, tag)
//
// All positional arguments are file paths.
func handleLess(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-b": true, "-h": true, "-j": true,
		"-x": true, "-y": true, "-z": true,
	}
	pathFlags := map[string]bool{
		"-k": true, "-o": true, "-O": true,
	}
	valueFlags := map[string]bool{
		"-p": true, "-t": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valueFlags, Placeholder: "<val>"},
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

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
