package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("sort", handleSort)
}

// longFlagValued lists long flags whose =value should be collapsed to =<val>.
var sortLongFlagsWithValue = map[string]bool{
	"--key": true, "--field-separator": true, "--buffer-size": true,
	"--output": true, "--temporary-directory": true, "--random-source": true,
	"--batch-size": true, "--compress-program": true,
}

// handleSort handles the sort command.
// Key flags (-k, --key) make the next arg <key>.
// Separator flags (-t, --field-separator) make the next arg <sep>.
// Size flags (-S, --buffer-size) make the next arg <size>.
// Path-consuming flags (-o, --output, -T, --temporary-directory, --random-source)
// make the next arg <path>.
// Batch size (--batch-size) makes the next arg N.
// Compress program (--compress-program) makes the next arg <cmd>.
// Long --flag=value forms collapse the value to <val>.
// All positionals are file paths.
func handleSort(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	keyFlags := map[string]bool{
		"-k": true, "--key": true,
	}
	sepFlags := map[string]bool{
		"-t": true, "--field-separator": true,
	}
	sizeFlags := map[string]bool{
		"-S": true, "--buffer-size": true,
	}
	pathFlags := map[string]bool{
		"-o": true, "--output": true,
		"-T": true, "--temporary-directory": true,
		"--random-source": true,
	}
	numericFlags := map[string]bool{
		"--batch-size": true,
	}
	cmdFlags := map[string]bool{
		"--compress-program": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: keyFlags, Placeholder: "<key>"},
		{Flags: sepFlags, Placeholder: "<sep>"},
		{Flags: sizeFlags, Placeholder: "<size>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: cmdFlags, Placeholder: "<cmd>"},
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

		// Handle --flag=value forms
		if strings.HasPrefix(tok, "--") {
			if eqIdx := strings.Index(tok, "="); eqIdx >= 0 {
				name := tok[:eqIdx]
				if sortLongFlagsWithValue[name] {
					result = append(result, name+"=<val>")
					i++
					continue
				}
			}
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

		// Positional: input file
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
