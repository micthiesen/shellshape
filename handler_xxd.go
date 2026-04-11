package shellshape

func init() {
	Register("xxd", handleXxd)
}

// handleXxd handles the xxd hex dump command.
// Numeric flags (-l, -s, -c, -g, -o and long forms) consume next arg as N.
// The -n/--name flag consumes next arg as <val>.
// The -R flag consumes next arg as <val> (color mode).
// All positionals are file paths (infile, outfile).
func handleXxd(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

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

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
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
