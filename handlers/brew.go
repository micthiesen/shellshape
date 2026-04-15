package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("brew", handleBrew, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleBrew handles brew (Homebrew) subcommand arguments.
// Since brew is in subcommandExecutables, the subcommand (install, uninstall,
// upgrade, search, info, list, tap, untap, services, etc.) is already consumed.
//
// For install/uninstall/upgrade/info/reinstall/deps/pin/unpin/fetch/audit and
// similar: positionals are formula or cask names, collapsed to <pkg>.
// For search: the positional is a query string, collapsed to <query>.
// For tap/untap: positionals are tap names (user/repo), kept verbatim.
// For services: first positional is the action (list/start/stop/...) kept
// verbatim; subsequent positionals are service names collapsed to <pkg>.
// Boolean flags (--cask, --formula, --HEAD, --force, --verbose, etc.) are kept.
// Very few flags consume values; those that do collapse the value.
func handleBrew(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Subcommands where positionals are formula/cask names, collapsed to <pkg>.
	pkgSubcommands := map[string]bool{
		"install": true, "uninstall": true, "remove": true, "rm": true,
		"upgrade": true, "info": true, "reinstall": true,
		"deps": true, "uses": true, "pin": true, "unpin": true,
		"edit": true, "home": true, "log": true, "fetch": true, "audit": true,
		"link": true, "unlink": true, "switch": true,
	}

	// Subcommands whose positionals are kept verbatim (tap names like
	// `user/repo`, or the bare `list` that still accepts formulas).
	verbatimSubcommands := map[string]bool{
		"tap": true, "untap": true, "list": true, "ls": true,
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
	isServicesSubcommand := subcommand == "services"

	var result []string
	seenServicesAction := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional handling depends on subcommand.
		if isSearchSubcommand {
			result = append(result, "<query>")
			i++
			continue
		}

		if isServicesSubcommand {
			if !seenServicesAction {
				// list/start/stop/restart/run/cleanup/info/kill — structural.
				result = append(result, tok)
				seenServicesAction = true
			} else {
				result = append(result, "<pkg>")
			}
			i++
			continue
		}

		if pkgSubcommands[subcommand] {
			result = append(result, "<pkg>")
			i++
			continue
		}

		if verbatimSubcommands[subcommand] {
			result = append(result, tok)
			i++
			continue
		}

		// Other subcommands: generic classification.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
