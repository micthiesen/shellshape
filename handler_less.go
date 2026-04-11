package shellshape

func init() {
	for _, name := range []string{"less", "more"} {
		Register(name, handleLess)
	}
}

// handleLess handles less and more.
// Most flags are boolean. A few flags consume the next token:
//   - Numeric: -b, -h, -j, -x, -y, -z (buffer, scroll, target, tabs, scroll limit, window)
//   - Path: -k, -o, -O (lesskey file, log file)
//   - Value: -p, -t (search pattern, tag)
//
// All positional arguments are file paths.
func handleLess(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	numericFlags := map[string]bool{
		"-b": true, "-h": true, "-j": true,
		"-x": true, "-y": true, "-z": true,
	}
	pathFlags := map[string]bool{
		"-k": true, "-o": true, "-O": true,
	}
	valueFlags := map[string]bool{
		"-p": true, "-t": true,
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

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if valueFlags[tok] {
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

		// Positional: file path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
