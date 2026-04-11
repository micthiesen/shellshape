package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("mtr", handleMtr)
}

// handleMtr handles mtr (network diagnostic combining traceroute and ping).
// Numeric flags (-c, -i, -s, -m, -f, -p, -z, -B, -Q, and long forms) collapse next arg to N.
// Value flags (-a, -I, -o, -F, and long forms) collapse next arg to <val>.
// The single positional argument (host/IP) becomes <host>.
func handleMtr(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-c": true, "--report-cycles": true,
		"-i": true, "--interval": true,
		"-s": true, "--psize": true,
		"-m": true, "--max-ttl": true,
		"-f": true, "--first-ttl": true,
		"-p": true, "--port": true,
		"-z": true, "--timeout": true,
		"-B": true, "--bitpattern": true,
		"-Q": true, "--tos": true,
		"-G": true, "--gracetime": true,
		"--mark": true,
	}

	valueFlags := map[string]bool{
		"-a": true, "--address": true,
		"-I": true, "--interface": true,
		"-o": true, "--order": true,
		"-F": true, "--filename": true,
	}

	booleanDigitFlags := map[string]bool{
		"-4": true, "-6": true,
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

		if booleanDigitFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: the host
		result = append(result, "<host>")
		i++
	}

	result = append(result, redirects...)
	return result
}
