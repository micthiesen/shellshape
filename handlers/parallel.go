package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("parallel", handleParallel)
}

var parallelFusedNumericRE = regexp.MustCompile(`^-([jnLP])(\d+)$`)

// handleParallel handles GNU parallel.
// Parallel's own flags are processed first. The first non-flag positional
// becomes the command name (preserved verbatim). Subsequent positionals before
// a ::: or :::: delimiter are command template tokens (preserved verbatim).
// ::: and :::: are structural delimiters. Arguments after them are data values
// collapsed to <arg>. Subshells are always kept verbatim.
func handleParallel(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a numeric value.
	numericFlags := map[string]bool{
		"-j": true, "--jobs": true,
		"-n": true, "--max-args": true,
		"-L": true, "--max-lines": true,
		"-P":        true,
		"--timeout": true, "--retries": true, "--delay": true,
	}

	// Flags whose next token is a string/path value.
	stringFlags := map[string]bool{
		"-S": true, "--sshlogin": true, "--sshloginfile": true,
		"--colsep": true, "--tagstring": true,
		"--results": true, "--header": true,
		"--rpl": true, "--block": true,
	}

	pathFlags := map[string]bool{
		"-a": true, "--arg-file": true,
		"--workdir": true, "--joblog": true,
	}

	var result []string
	commandFound := false
	inArgSource := false  // after ::: (inline args)
	inFileSource := false // after :::: (arg files)

	i := 0
	for i < len(args) {
		tok := args[i]

		// ::: or :::: delimiter: switch to appropriate source mode.
		if tok == ":::" {
			result = append(result, tok)
			inArgSource = true
			inFileSource = false
			i++
			continue
		}
		if tok == "::::" {
			result = append(result, tok)
			inFileSource = true
			inArgSource = false
			i++
			continue
		}

		// In arg-source mode: data values get collapsed, subshells kept.
		if inArgSource {
			if shellshape.IsSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		// In file-source mode: tokens are file paths.
		if inFileSource {
			if shellshape.IsSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<path>")
			}
			i++
			continue
		}

		// Once command is found, remaining tokens before ::: are template tokens.
		if commandFound {
			if shellshape.IsSubshellToken(tok) {
				result = append(result, tok)
			} else {
				// Command template tokens (like {}, {.}.png) are structural.
				result = append(result, tok)
			}
			i++
			continue
		}

		// Still processing parallel's own flags.
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if stringFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<str>")
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// Fused numeric: -j4, -n2, -L1, -P8
		if m := parallelFusedNumericRE.FindStringSubmatch(tok); m != nil {
			result = append(result, "-"+m[1], "N")
			i++
			continue
		}

		// -0 is a boolean flag but isFlagToken doesn't catch it.
		if tok == "-0" {
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First non-flag positional is the command name.
		result = append(result, tok)
		commandFound = true
		i++
	}

	result = append(result, redirects...)
	return result
}
