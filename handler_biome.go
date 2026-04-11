package shellshape

func init() {
	Register("biome", handleBiome)
}

// handleBiome handles biome lint/format/check/ci/migrate.
// Tokens arrive after the subcommand has already been consumed.
// --config-path takes a path argument.
// --colors, --log-level, --diagnostic-level take a value argument.
// --write, --unsafe are boolean flags.
// All positionals are file/directory paths.
func handleBiome(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valueFlags := map[string]bool{
		"--config-path":      true,
		"--colors":           true,
		"--log-level":        true,
		"--diagnostic-level": true,
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

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file/directory path
		if tok == "." || tok == ".." {
			result = append(result, "<path>")
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
