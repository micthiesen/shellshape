package shellshape

func init() {
	Register("fnm", handleFnm, HandlerOptions{HasSubcommands: true})
}

// handleFnm handles fnm (Fast Node Manager) subcommand arguments.
// Since fnm is in subcommandExecutables, the subcommand (install, use, default,
// alias, etc.) is already consumed before this handler is called.
//
// For install/use/default/uninstall: the version positional collapses to <version>.
// For alias: first positional is <version>, second is kept verbatim (alias name).
// For exec: --using flag consumes a <version>, remaining positionals are the
// command to run and stay verbatim.
// For list/ls/current/env/completions/unalias/list-remote: positionals stay verbatim.
//
// Value flags (--node-dist-mirror, --log-level, --arch, --version-file-strategy,
// --progress, --resolve-engines) collapse their argument to <val>.
// Path flags (--fnm-dir) collapse to <path>.
// Shell flag (--shell) stays verbatim (structural).
// Boolean flags (--lts, --corepack-enabled, --latest, --json, --use) are kept as-is.
func handleFnm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Subcommands where the first positional is a version.
	versionSubcommands := map[string]bool{
		"install": true, "i": true,
		"use":       true,
		"default":   true,
		"uninstall": true, "uni": true,
		"alias": true,
	}

	// Flags whose next token collapses to <val>.
	valFlags := map[string]bool{
		"--node-dist-mirror":      true,
		"--log-level":             true,
		"--arch":                  true,
		"--version-file-strategy": true,
		"--progress":              true,
		"--resolve-engines":       true,
	}

	// Flags whose next token collapses to <path>.
	pathFlags := map[string]bool{
		"--fnm-dir": true,
	}

	// Flags whose next token stays verbatim (structural).
	verbatimFlags := map[string]bool{
		"--shell": true,
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

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// For exec subcommand, --using consumes a version.
		if tok == "--using" && subcommand == "exec" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<version>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional argument.
		positionalCount++

		if versionSubcommands[subcommand] {
			if subcommand == "alias" {
				if positionalCount == 1 {
					// First positional: version
					result = append(result, "<version>")
				} else {
					// Second positional: alias name, keep verbatim
					result = append(result, tok)
				}
			} else {
				// install/use/default/uninstall: version
				result = append(result, "<version>")
			}
		} else {
			// Other subcommands: keep verbatim
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
