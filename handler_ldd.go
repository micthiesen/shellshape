package shellshape

func init() {
	Register("ldd", handleLdd)
	Register("otool", handleOtool)
}

// handleLdd handles the ldd command (display shared library dependencies).
// All flags are boolean. All positional arguments are file paths.
func handleLdd(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: binary/library path — always a file
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleOtool handles the otool command (display object file info on macOS).
// -arch consumes one argument (architecture name).
// -p consumes one argument (symbol name).
// -s consumes two arguments (segment and section names).
// All other flags are boolean. All positionals are file paths.
func handleOtool(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if tok == "-arch" || tok == "--arch" {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if tok == "-p" {
			result, i = consumeFlagArg(tok, args, i, result, "<sym>")
			continue
		}

		if tok == "-s" {
			// -s takes two arguments: segname sectname
			result = append(result, tok)
			i++
			for n := 0; n < 2 && i < len(args); n++ {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: object file path — always a file
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
