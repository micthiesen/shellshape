package shellshape

// handleBrew handles brew (Homebrew) subcommand arguments.
// Since brew is in subcommandExecutables, the subcommand (install, uninstall,
// upgrade, search, info, list, tap, untap, services, etc.) is already consumed.
//
// For install/uninstall/upgrade/info/reinstall/deps/pin/unpin and similar:
// positionals are formula or cask names, kept verbatim (structural).
// For search: the positional is a query string, collapsed to <query>.
// For tap/untap: positionals are tap names (user/repo), kept verbatim.
// Boolean flags (--cask, --formula, --HEAD, --force, --verbose, etc.) are kept.
// Very few flags consume values; those that do collapse the value.
func handleBrew(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Subcommands where positionals are formula/cask names (structural).
	formulaSubcommands := map[string]bool{
		"install": true, "uninstall": true, "remove": true, "rm": true,
		"upgrade": true, "info": true, "reinstall": true,
		"deps": true, "uses": true, "pin": true, "unpin": true,
		"edit": true, "home": true, "log": true, "fetch": true, "audit": true,
		"list": true, "ls": true,
		"link": true, "unlink": true, "switch": true,
		"tap": true, "untap": true,
		"services": true,
	}

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"--appdir":  true,
		"--fontdir": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"--language": true,
	}

	isSearchSubcommand := subcommand == "search"

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
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

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional handling depends on subcommand.
		if isSearchSubcommand {
			result = append(result, "<query>")
			i++
			continue
		}

		if formulaSubcommands[subcommand] {
			// Formula/cask/tap names and service names are structural.
			result = append(result, tok)
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
