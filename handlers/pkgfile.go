package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("pkgfile", handlePkgfile)
}

// handlePkgfile handles pkgfile (Arch Linux package file search tool).
//
// Boolean flags: -u (update), -s (search), -r (regex), -b (binaries),
// -d (directories), -v (verbose), -q (quiet), -g (glob).
// Flags consuming next token: -l/--list (package name → <val>),
// -R/--repo (repo name → <val>).
// Positional argument is a filename/pattern → <pattern>.
func handlePkgfile(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-l":     true,
		"--list": true,
		"-R":     true,
		"--repo": true,
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

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: filename or pattern to search for
		result = append(result, "<pattern>")
		i++
	}

	result = append(result, redirects...)
	return result
}
