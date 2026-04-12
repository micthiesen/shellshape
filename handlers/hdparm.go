package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("hdparm", handleHdparm)
}

// handleHdparm handles hdparm HDD parameter commands.
// Numeric flags: -S, -B, -M, -a, -m, -W, -A, -c, -d, -u, -X, -p, -r
// These flags optionally take a numeric argument (get vs set mode).
// If the next token is a number, consume it as N; otherwise the flag is boolean.
// Value flags: --security-set-pass, --security-erase, etc. → <val>
// Boolean flags: -I, -t, -T, -C, -g, -i, -y, -Y, -Z, etc.
// Positional is always a device path → <path>.
func handleHdparm(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags that optionally take a numeric argument
	numericOptFlags := map[string]bool{
		"-S": true, "-B": true, "-M": true,
		"-a": true, "-m": true, "-W": true,
		"-A": true, "-c": true, "-d": true,
		"-u": true, "-X": true, "-p": true,
		"-r": true,
	}

	// Long flags that take a value argument
	valFlags := map[string]bool{
		"--security-set-pass":       true,
		"--security-erase":          true,
		"--security-erase-enhanced": true,
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

		// Numeric-optional flags: consume next token only if it's a number
		if numericOptFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) && shellshape.NumberRE.MatchString(args[i]) {
				result = append(result, "N")
				i++
			}
			continue
		}

		// Long value flags
		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path
		if shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
