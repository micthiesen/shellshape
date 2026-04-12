package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("snapper", handleSnapper, shellshape.HandlerOptions{HasSubcommands: true})
}

var snapperIDRE = regexp.MustCompile(`^\d+([.]{2,3}\d+|-\d+)?$`)

// handleSnapper handles snapper (Btrfs snapshot manager).
//
// Subcommands: list, create, delete, diff, undochange, status, rollback,
// get-config, set-config, list-configs, cleanup, modify, setup-quota.
//
// Global flags that may appear before the subcommand:
//
//	-c/--config <name> (config name → <val>)
//	-v/--verbose, -q/--quiet, --no-dbus, --csvout, --jsonout (boolean)
//
// Snapshot IDs (plain numbers, N..M ranges, N-M ranges) → N.
// -t/--type values (pre/post/single) are kept verbatim.
// -d/--description, -u/--userdata, --command, --cleanup-algorithm → <val>.
// Paths → classifyToken.
func handleSnapper(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose values are kept verbatim (small finite set).
	verbatimFlags := map[string]bool{
		"-t": true, "--type": true,
	}

	valFlags := map[string]bool{
		"-c": true, "--config": true,
		"-d": true, "--description": true,
		"-u": true, "--userdata": true,
		"--command":           true,
		"--cleanup-algorithm": true,
		"--from":              true,
	}

	numericFlags := map[string]bool{
		"--pre-num": true,
		"--num":     true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	i := 0

	// Handle global flags consumed as subcommand by the normalizer.
	if shellshape.IsFlagToken(subcommand) {
		if placeholder, ok := shellshape.MatchFlagCategory(subcommand, categories); ok {
			if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
				result = append(result, placeholder)
				i++
			}
		} else if verbatimFlags[subcommand] {
			if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
				result = append(result, args[i])
				i++
			}
		}
		// Emit real subcommand.
		if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
			result = append(result, args[i])
			i++
		}
	}

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) && !shellshape.IsFlagToken(args[i]) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: snapshot IDs (numbers, ranges) → N, else → <val>.
		if snapperIDRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
