package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("fd", handleFd)
}

// handleFd handles fd (fast find).
// First positional → <pattern>, second positional → path (ClassifyToken).
// -t/--type, -e/--extension → <val>.
// -E/--exclude → <pattern>.
// -d/--max-depth, --min-depth → N.
// -x/--exec, -X/--exec-batch → consumes rest of args as <cmd...>.
// Boolean: -H, -I, -u, -l, -a, -L, -p, -g, etc.
func handleFd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-t": true, "--type": true,
		"-e": true, "--extension": true,
	}
	patternFlags := map[string]bool{
		"-E": true, "--exclude": true,
	}
	numericFlags := map[string]bool{
		"-d": true, "--max-depth": true, "--min-depth": true,
	}
	execFlags := map[string]bool{
		"-x": true, "--exec": true, "-X": true, "--exec-batch": true,
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

		// Exec flags consume the rest of args
		if execFlags[tok] {
			result = append(result, tok, "<cmd...>")
			// Skip all remaining args
			i = len(args)
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if patternFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<pattern>")
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
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
