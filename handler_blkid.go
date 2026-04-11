package shellshape

func init() {
	Register("blkid", handleBlkid)
}

// handleBlkid handles the blkid block device identifier command.
// Flags that take value arguments (-L, -U, -n, -o, -s, -t, -u, -H and long forms)
// collapse their argument to <val>.
// Cache file flags (-c, --cache-file) collapse to <path>.
// Numeric flags (-O, -S and long forms) collapse to N.
// Remaining positionals (device paths) use classifyToken.
func handleBlkid(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	valFlags := map[string]bool{
		"-L": true, "--label": true,
		"-U": true, "--uuid": true,
		"-n": true, "--match-types": true,
		"-o": true, "--output": true,
		"-s": true, "--match-tag": true,
		"-t": true, "--match-token": true,
		"-u": true, "--usages": true,
		"-H": true, "--hint": true,
	}
	pathFlags := map[string]bool{
		"-c": true, "--cache-file": true,
	}
	numericFlags := map[string]bool{
		"-O": true, "--offset": true,
		"-S": true, "--size": true,
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

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
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

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional (device path)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
