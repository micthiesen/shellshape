package shellshape

func init() {
	Register("sst", handleSst, HandlerOptions{HasSubcommands: true})
}

// handleSst handles sst (SST Ion/v3) infrastructure commands.
// The subcommand (dev, deploy, build, remove, secret, shell, tunnel, etc.)
// is already consumed by the normalizer since sst is in subcommandExecutables.
//
// --stage value is kept verbatim (stage names like "prod", "dev" are structural).
// --profile and --region consume the next token and collapse it to <val>.
// --verbose and other boolean flags are kept as-is.
// First two positionals stay verbatim (covers sub-subcommands and identifiers
// like secret names). Third+ positionals collapse to <str>.
func handleSst(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token should be kept verbatim.
	verbatimValueFlags := map[string]bool{
		"--stage": true,
	}

	// Flags whose next token should collapse to <val>.
	collapseValueFlags := map[string]bool{
		"--profile": true,
		"--region":  true,
	}

	var result []string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		if verbatimValueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if collapseValueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional argument.
		positionalCount++
		if positionalCount <= 2 {
			// First two positionals: classify normally (sub-subcommands and
			// identifiers like secret names stay verbatim, paths collapse).
			result = append(result, classifyToken(tok))
		} else {
			// Third+ positionals collapse to <str> (e.g. secret values).
			result = append(result, "<str>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
