package shellshape

func init() {
	for _, name := range []string{"ping", "ping6"} {
		Register(name, handlePing)
	}
}

// handlePing handles ping and ping6.
// Numeric flags (-c, -G, -g, -h, -i, -l, -m, -s, -t, -T, -W, -z) collapse next arg to N.
// Value flags (-b, -I, -k, -K, -M, -P, -p, -S) collapse next arg to <val>.
// The single positional argument (host/IP) becomes <host>.
func handlePing(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-c": true, "-G": true, "-g": true, "-h": true,
		"-i": true, "-l": true, "-m": true, "-s": true,
		"-t": true, "-T": true, "-W": true, "-z": true,
	}

	valueFlags := map[string]bool{
		"-b": true, "-I": true, "-k": true, "-K": true,
		"-M": true, "-P": true, "-p": true, "-S": true,
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
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if valueFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
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
