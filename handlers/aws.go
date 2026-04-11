package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("aws", handleAws, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleAws handles aws CLI commands.
// Tokens arrive after the first subcommand (service name like s3, lambda, ec2)
// since aws is in subcommandExecutables. The first positional is kept verbatim
// as the second subcommand. Global flags like --region and --output keep their
// values verbatim. --query collapses to <query>. Most other flags with values
// collapse to <val>. Positionals use classifyToken.
func handleAws(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose values are kept verbatim (structural, small set of values).
	verbatimValueFlags := map[string]bool{
		"--region": true,
		"--output": true,
		"--color":  true,
	}

	// Boolean flags (no argument consumed).
	booleanFlags := map[string]bool{
		"--debug":            true,
		"--no-verify-ssl":    true,
		"--no-paginate":      true,
		"--no-sign-request":  true,
		"--recursive":        true,
		"--dryrun":           true,
		"--dry-run":          true,
		"--delete":           true,
		"--only-show-errors": true,
		"--no-include-email": true,
		"--force":            true,
		"--exact-timestamps": true,
		"--quiet":            true,
		"--version":          true,
	}

	// Flags whose values collapse to a specific placeholder.
	queryFlags := map[string]bool{
		"--query": true,
	}

	// Flags whose values are numeric.
	numericFlags := map[string]bool{
		"--cli-read-timeout":    true,
		"--cli-connect-timeout": true,
		"--page-size":           true,
		"--max-items":           true,
		"--starting-token":      true,
	}

	// Flags whose values are paths.
	pathFlags := map[string]bool{
		"--ca-bundle": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: queryFlags, Placeholder: "<query>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	var result []string
	subcmdAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !subcmdAssigned {
				subcmdAssigned = true
			}
			i++
			continue
		}

		// Boolean Flags: keep verbatim, no value consumed.
		if booleanFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Verbatim value Flags: keep flag and value as-is.
		if verbatimValueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Other long flags with values: collapse to <val>.
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			// Long flags (--foo) that aren't boolean consume next token as value.
			if len(tok) > 2 && tok[:2] == "--" && !booleanFlags[tok] {
				i++
				if i < len(args) && !shellshape.IsFlagToken(args[i]) {
					if shellshape.IsSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "<val>")
					}
				} else {
					// Next token is a flag or missing; don't consume.
					continue
				}
			}
			i++
			continue
		}

		// Positional: first one is the second subcommand, kept verbatim.
		if !subcmdAssigned {
			result = append(result, tok)
			subcmdAssigned = true
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
