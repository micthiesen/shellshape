package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("lsblk", handleLsblk)
}

// handleLsblk handles the lsblk (list block devices) command.
// -o/--output consume a column list (<columns>).
// -e/--exclude and -I/--include consume major device numbers (N for single, <columns> for csv).
// -x/--sort consume a column name (<column>).
// -w/--width consume a numeric width (N).
// All other flags are boolean. Positionals are device paths via classifyToken.
func handleLsblk(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	columnListFlags := map[string]bool{
		"-o": true, "--output": true,
	}
	majorListFlags := map[string]bool{
		"-e": true, "--exclude": true,
		"-I": true, "--include": true,
	}
	sortFlags := map[string]bool{
		"-x": true, "--sort": true,
	}
	numericFlags := map[string]bool{
		"-w": true, "--width": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: columnListFlags, Placeholder: "<columns>"},
		{Flags: sortFlags, Placeholder: "<column>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if majorListFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				val := args[i]
				if strings.Contains(val, ",") {
					result = append(result, "<columns>")
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: device path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
