package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("diff", handleDiff)
}

var fusedDiffNumericRE = regexp.MustCompile(`^-([CUW])(\d+)$`)

// handleDiff handles the diff command.
// Numeric flags (-C, -U, -W, --context, --unified, --tabsize) collapse next arg to N.
// Pattern flags (-I, -F, -x) collapse next arg to <pattern>.
// Label/string flags (-L, -D, --changed-group-format) collapse next arg to <str>.
// File flags (-X, -S) collapse next arg to <path>.
// Algorithm flags (-A, --algorithm) keep next arg verbatim.
// Fused numeric forms like -C3 are split into -C and N.
// Positionals use classifyToken (typically file/directory paths).
func handleDiff(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-C": true, "-U": true, "-W": true,
		"--context": true, "--unified": true, "--tabsize": true, "--width": true,
	}
	patternFlags := map[string]bool{
		"-I": true, "--ignore-matching-lines": true,
		"-F": true, "--show-function-line": true,
		"-x": true, "--exclude": true,
	}
	stringFlags := map[string]bool{
		"-L": true, "--label": true,
		"-D": true, "--ifdef": true,
		"--changed-group-format": true,
	}
	fileFlags := map[string]bool{
		"-X": true, "--exclude-from": true,
		"-S": true, "--starting-file": true,
	}
	algorithmFlags := map[string]bool{
		"-A": true, "--algorithm": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: patternFlags, Placeholder: "<pattern>"},
		{Flags: stringFlags, Placeholder: "<str>"},
		{Flags: fileFlags, Placeholder: "<path>"},
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

		if algorithmFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Fused numeric: -C3, -U5, -W120
		if m := fusedDiffNumericRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "N")
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional (file/dir path, or - for stdin)
		if tok == "-" {
			result = append(result, "-")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
