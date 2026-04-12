package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("rg", handleRg)
}

var rgFusedNumericRE = regexp.MustCompile(`^-([ABC])(\d+)$`)

// handleRg handles ripgrep (rg).
// First positional is the search pattern → <pattern>.
// Subsequent positionals are paths → ClassifyToken.
// -e <pattern>, -g/--glob <pattern> → <pattern>.
// -t/--type → <val> (structural).
// -A/-B/-C N, --max-count N → N. Fused -A3 etc.
func handleRg(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	patternFlags := map[string]bool{"-e": true, "-g": true, "--glob": true}
	valFlags := map[string]bool{"-t": true, "--type": true}
	numericFlags := map[string]bool{
		"-A": true, "-B": true, "-C": true,
		"--max-count": true, "--after-context": true,
		"--before-context": true, "--context": true,
	}

	var result []string
	patternAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			patternAssigned = true
			i++
			continue
		}

		if patternFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<pattern>")
			if tok == "-e" {
				patternAssigned = true
			}
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// Fused numeric: -A3, -B10, -C2
		if m := rgFusedNumericRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "N")
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !patternAssigned {
			result = append(result, "<pattern>")
			patternAssigned = true
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
