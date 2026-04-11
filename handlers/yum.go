package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"yum", "dnf"} {
		shellshape.Register(name, handleYum, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

// handleYum handles yum and dnf subcommand arguments.
// yum and dnf are RPM-based package managers with nearly identical grammars.
// Since HasSubcommands is set, the subcommand (install, remove, search, etc.)
// is already consumed before this handler is called.
//
// Package names are structural (kept verbatim) for install, remove, update, etc.
// Search queries collapse to <query>. Flags that consume values collapse their
// argument appropriately. Boolean flags are kept verbatim.
func handleYum(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Subcommands where positionals are package names (structural).
	packageSubcommands := map[string]bool{
		"install": true, "reinstall": true,
		"remove": true, "erase": true, "autoremove": true,
		"update": true, "upgrade": true, "downgrade": true,
		"info": true, "deplist": true,
		"mark": true, "swap": true,
		"download": true, "builddep": true,
		// list/clean/repolist/history: positionals are keywords, kept verbatim
		"list": true, "clean": true,
		"repolist": true, "repoinfo": true,
		"history": true, "group": true,
		"makecache": true, "check-update": true,
		"module": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"--enablerepo":   true,
		"--disablerepo":  true,
		"--exclude":      true,
		"--releasever":   true,
		"--setopt":       true,
		"--repo":         true,
		"--repoid":       true,
		"--advisory":     true,
		"--cve":          true,
		"--sec-severity": true,
	}

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"-c": true, "--config": true,
		"--installroot": true,
		"--downloaddir": true,
		"--destdir":     true,
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

		if packageSubcommands[subcommand] {
			// Package names and subcommand keywords are structural, but
			// pure numbers (e.g. transaction IDs in history) should collapse.
			if shellshape.NumberRE.MatchString(tok) {
				result = append(result, "N")
			} else {
				result = append(result, tok)
			}
			i++
			continue
		}

		// Other/unknown subcommands: generic classification.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
