package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("vagrant", handleVagrant, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleVagrant handles the vagrant CLI.
// vagrant has subcommands (up, ssh, box, snapshot, plugin, etc.) which are
// already extracted before this handler is called.
//
// Some subcommands (box, snapshot, plugin) have a second-level subcommand
// (box add, snapshot save, plugin install) which is kept verbatim.
// After the second-level subcommand, positionals are data and collapse to <val>.
//
// For init, positionals (box name, optional URL) are data and collapse.
// For other subcommands (up, halt, ssh, destroy, etc.), the positional is
// typically a VM name from the Vagrantfile and is kept verbatim.
func handleVagrant(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"--provider":       true,
		"--provision-with": true,
		"-c":               true, "--command": true,
		"--checksum":      true,
		"--checksum-type": true,
		"--name":          true,
		"--box-version":   true,
		"--cacert":        true,
		"--capath":        true,
		"--cert":          true,
		"--box":           true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"--output": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	// Subcommands that have a second-level subcommand.
	hasSecondSub := map[string]bool{
		"box": true, "snapshot": true, "plugin": true,
	}

	// Subcommands where all positionals are data (collapse to <val>/classifyToken).
	allPositionalsData := map[string]bool{
		"init": true,
	}

	var result []string
	secondSubTaken := !hasSecondSub[subcommand]
	collapsePositionals := allPositionalsData[subcommand] || hasSecondSub[subcommand]

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Second-level subcommand: keep verbatim.
		if !secondSubTaken {
			result = append(result, tok)
			secondSubTaken = true
			i++
			continue
		}

		// Positional argument.
		if collapsePositionals {
			if shellshape.IsSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<val>")
			}
		} else {
			// VM name: keep verbatim.
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
