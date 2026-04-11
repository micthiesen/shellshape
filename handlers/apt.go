package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"apt", "apt-get", "apt-cache"} {
		shellshape.Register(name, handleApt, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

// handleApt handles apt, apt-get, and apt-cache subcommand arguments.
// Since these are registered with HasSubcommands, the subcommand (install,
// remove, update, search, show, etc.) is already consumed.
//
// For install/remove/purge/reinstall/download/show/depends/rdepends/policy/showpkg:
// positionals are package names, kept verbatim (structural).
// For search: the positional is a query string, collapsed to <query>.
// For update/upgrade/dist-upgrade/full-upgrade/autoremove/autoclean/clean/check:
// no positionals expected, just boolean flags.
// Flags -t/--target-release and -o/--option consume the next token as <val>.
func handleApt(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Subcommands where positionals are package names (structural).
	packageSubcommands := map[string]bool{
		"install": true, "remove": true, "purge": true, "reinstall": true,
		"download": true, "changelog": true,
		"show": true, "showpkg": true, "showsrc": true,
		"depends": true, "rdepends": true,
		"policy": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"-t": true, "--target-release": true, "--default-release": true,
		"-o": true, "--option": true,
	}

	isSearchSubcommand := subcommand == "search"

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
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

		if packageSubcommands[subcommand] {
			// Package names are structural.
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
