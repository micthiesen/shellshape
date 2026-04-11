package shellshape

// handleGh handles the gh (GitHub CLI) command.
// gh is in subcommandExecutables, so the first subcommand (pr, issue, repo, api)
// is already extracted. This handler receives tokens after that first subcommand.
//
// Many gh commands have a second subcommand (pr create, issue list, repo clone)
// which is kept verbatim as structural.
//
// String flags (--title, --body) collapse to <str>.
// Value flags (--repo, --assignee, --label, --json, --jq, etc.) collapse to <val>.
// Numeric flags (--limit) collapse to N.
// Boolean flags (--web, --draft, --fill, etc.) are kept verbatim.
// Positional PR/issue numbers become N. URLs go through classifyToken.
func handleGh(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	strFlags := map[string]bool{
		"--title": true, "-t": true,
		"--body": true, "-b": true,
	}

	valFlags := map[string]bool{
		"--repo": true, "-R": true,
		"--assignee": true, "-a": true,
		"--label": true, "-l": true,
		"--json": true,
		"--jq": true, "-q": true,
		"--base": true, "-B": true,
		"--head": true, "-H": true,
		"--milestone": true, "-m": true,
		"--reviewer": true, "-r": true,
		"--project": true, "-p": true,
		"--template": true, "-T": true,
		"--body-file": true, "-F": true,
		"--state": true, "-s": true,
		"--author": true,
		"--search": true,
		"--sort": true,
		"--order": true,
		"--hostname": true,
		"-f": true, "--raw-field": true,
		"--field": true,
		"-X": true, "--method": true,
		"--header": true,
		"--recover": true,
		"--app": true,
		"--color": true,
		"--language": true,
		"--topic": true,
		"--visibility": true,
	}

	numericFlags := map[string]bool{
		"--limit": true, "-L": true,
	}

	// Subcommands where the first positional is NOT a second subcommand.
	// For these, all positionals are data and should be classified.
	noSecondSubcmd := map[string]bool{
		"api": true, "status": true, "completion": true,
	}

	var result []string
	subcommandTaken := noSecondSubcmd[subcommand]

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if strFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<str>")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first positional is the second subcommand, kept verbatim
		if !subcommandTaken {
			result = append(result, tok)
			subcommandTaken = true
			i++
			continue
		}

		// Remaining positionals: classify (numbers become N, URLs get classified, etc.)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
