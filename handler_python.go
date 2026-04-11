package shellshape

func init() {
	for _, name := range []string{"python", "python3"} {
		Register(name, handlePython)
	}
}

// handlePython handles python and python3.
// -c takes inline code. -m takes a module name. -W and -X take a value.
// The first positional is the script path; everything after the script
// (or after -c <code> / -m <module>) is a script argument collapsed to <arg>.
func handlePython(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-W":                      true,
		"-X":                      true,
		"--check-hash-based-pycs": true,
	}

	var result []string
	// "script" means a script file was seen; args after it are opaque <arg>.
	// "module" means -m <mod> was seen; args after it get classifyToken.
	scriptSeen := false
	moduleSeen := false
	doubleDash := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if tok == "--" {
			result = append(result, "--")
			doubleDash = true
			i++
			continue
		}

		if doubleDash || scriptSeen {
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		if moduleSeen {
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, classifyToken(tok))
			}
			i++
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			scriptSeen = true
			i++
			continue
		}

		// -c: inline code, marks script seen (everything after is <arg>).
		if tok == "-c" {
			result, i = consumeFlagArg(tok, args, i, result, "<code>")
			scriptSeen = true
			continue
		}

		// -m: module name is kept verbatim (it's structural, like a subcommand).
		// Everything after the module name gets classifyToken.
		if tok == "-m" {
			result = append(result, "-m")
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			moduleSeen = true
			continue
		}

		// Value-consuming flags.
		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Any other flag is boolean.
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional: the script path.
		result = append(result, classifyToken(tok))
		scriptSeen = true
		i++
	}

	result = append(result, redirects...)
	return result
}
