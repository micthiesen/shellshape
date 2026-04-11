package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("xxd", handleXxd)
}

// handleXxd handles the xxd hex dump command.
// Numeric flags (-l, -s, -c, -g, -o and long forms) consume next arg as N.
// The -n/--name flag consumes next arg as <val>.
// The -R flag consumes next arg as <val> (color mode).
// All positionals are file paths (infile, outfile).
func handleXxd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-l": true, "--len": true,
		"-s": true, "--seek": true,
		"-c": true, "--cols": true,
		"-g": true, "--groupsize": true,
		"-o": true,
	}

	valFlags := map[string]bool{
		"-n": true, "--name": true,
		"-R": true,
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

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a file path (infile or outfile)
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
