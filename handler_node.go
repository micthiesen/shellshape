package shellshape

func init() {
	Register("node", handleNode)
}

// handleNode handles node (and compatible runtimes).
// -e/-c take inline code. -r/--require, --loader, --env-file, --import take
// paths. -C/--conditions take a value. The first positional is the script
// path; everything after the script (or after -e/-c <code>) is a script
// argument collapsed to <arg>.
func handleNode(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-r": true, "--require": true,
		"--loader": true, "--experimental-loader": true,
		"--env-file": true, "--env-file-if-exists": true,
		"--import":       true,
		"--cpu-prof-dir": true, "--cpu-prof-name": true,
		"--heap-prof-dir": true, "--heap-prof-name": true,
		"--diagnostic-dir":           true,
		"--experimental-sea-config":  true,
		"--experimental-config-file": true,
		"--build-snapshot-config":    true,
		"--localstorage-file":        true,
		"--icu-data-dir":             true,
	}

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-C": true, "--conditions": true,
		"--input-type":                        true,
		"--experimental-specifier-resolution": true,
		"--inspect-publish-uid":               true,
		"--dns-result-order":                  true,
		"--disable-proto":                     true,
		"--disable-warning":                   true,
		"--heapsnapshot-signal":               true,
	}

	// Flags whose next token is inline code.
	codeFlags := map[string]bool{
		"-e": true, "--eval": true,
		"-c": true, "--check": true,
	}

	var result []string
	scriptSeen := false // true after the script path or -e/-c <code>
	doubleDash := false

	i := 0
	for i < len(args) {
		tok := args[i]

		// After --, everything is a script arg.
		if tok == "--" {
			result = append(result, "--")
			scriptSeen = true
			i++
			continue
		}

		if doubleDash || scriptSeen {
			// Script arguments: collapse all to <arg>.
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<arg>")
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

		// Code flags: -e/-c consume the next token as <code>.
		if codeFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<code>")
				scriptSeen = true
				i++
			}
			continue
		}

		// Path-consuming flags.
		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		// Value-consuming flags.
		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		// Any other flag (boolean or --flag=value handled by classifyToken).
		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
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
