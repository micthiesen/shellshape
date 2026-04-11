package shellshape

func init() {
	for _, name := range []string{"less", "more"} {
		Register(name, handleLess)
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
	args, redirects := splitRedirects(tokens)

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

	categories := []flagCategory{
		{numericFlags, "N"},
		{pathFlags, "<path>"},
		{valueFlags, "<val>"},
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
