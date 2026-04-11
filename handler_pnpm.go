package shellshape

func init() {
	for _, name := range []string{"pnpm", "npm", "yarn"} {
		Register(name, handlePnpm, HandlerOptions{HasSubcommands: true})
	}
}

// handlePnpm handles pnpm, npm, and yarn commands.
// The subcommand (install, add, remove, run, exec, dlx, etc.) is already
// extracted by normalizeSingleCommand and passed via the subcommand parameter.
//
// For add/remove: all positionals collapse to <pkg>.
// For run: first positional is the script name (kept verbatim), rest pass through.
// For exec/dlx: first positional is the command (kept verbatim), rest collapse to <arg>.
// For other subcommands: positionals use classifyToken.
//
// Key flags: --filter/-F (value), --config (path), -C (path),
// --recursive/-r, --frozen-lockfile, --prod/--production, --save-dev/-D,
// --save-exact/-E (all boolean).
func handlePnpm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags that consume the next token as a value.
	valFlags := map[string]bool{
		"--filter": true, "-F": true,
		"--package": true,
	}

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"--config": true,
		"-C":       true,
	}

	// Subcommands where all positionals are package names.
	pkgSubcommands := map[string]bool{
		"add": true, "remove": true, "rm": true, "uninstall": true, "un": true,
		"install": true, "i": true,
	}

	// Subcommands where first positional is a command name (verbatim)
	// and the rest passes through verbatim (arguments to the command).
	execSubcommands := map[string]bool{
		"exec": true, "run": true,
	}

	// Subcommands where first positional is a command name (verbatim)
	// and remaining positionals collapse to <arg>.
	dlxSubcommands := map[string]bool{
		"dlx": true,
	}

	var result []string
	firstPositional := true
	passthrough := false // after -- or after first positional in run/exec

	i := 0
	for i < len(args) {
		tok := args[i]

		// -- separator: pass through everything after it verbatim
		if tok == "--" {
			result = append(result, args[i:]...)
			break
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			if firstPositional {
				firstPositional = false
			}
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

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional handling depends on subcommand.
		if pkgSubcommands[subcommand] {
			result = append(result, "<pkg>")
			i++
			continue
		}

		if execSubcommands[subcommand] {
			if firstPositional {
				// Script/command name: keep verbatim.
				result = append(result, tok)
				firstPositional = false
				passthrough = true
			} else if passthrough {
				// After the name in run/exec, pass through verbatim.
				result = append(result, args[i:]...)
				break
			}
			i++
			continue
		}

		if dlxSubcommands[subcommand] {
			if firstPositional {
				// Command name: keep verbatim.
				result = append(result, tok)
				firstPositional = false
			} else {
				// Remaining args are data (project name, etc.).
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		// Other subcommands: generic classification.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
