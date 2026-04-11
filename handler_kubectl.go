package shellshape

// handleKubectl handles kubectl subcommand arguments.
// Since kubectl is in subcommandExecutables, the subcommand (get, describe,
// apply, delete, logs, exec, etc.) is already consumed before this handler.
//
// For subcommands that take a resource type (get, describe, delete, edit, etc.),
// the first non-flag positional is the resource type, kept verbatim.
// For subcommands that take a resource name directly (logs, exec, attach, etc.),
// there is no resource type.
// Remaining positionals are resource names, collapsed to <word>.
//
// Key flags: -n/--namespace (value), -l/--selector (selector), -o/--output
// (verbatim), -f/--filename (path), --context (value), --kubeconfig (path),
// -c/--container (value).
func handleKubectl(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Subcommands where the first positional is a resource name (no resource type).
	nameFirstSubcmds := map[string]bool{
		"logs":         true,
		"exec":         true,
		"attach":       true,
		"port-forward": true,
		"cp":           true,
		"debug":        true,
		"run":          true,
	}

	// Flags whose values are kept verbatim (small finite set of values).
	verbatimValueFlags := map[string]bool{
		"-o": true, "--output": true,
	}

	// Flags whose values collapse to <selector>.
	selectorFlags := map[string]bool{
		"-l": true, "--selector": true,
		"--field-selector": true,
	}

	// Flags whose values collapse to <val>.
	valFlags := map[string]bool{
		"-n": true, "--namespace": true,
		"-c": true, "--container": true,
		"--context":  true,
		"--sort-by":  true,
		"--template": true, "--go-template": true,
		"--type":     true,
		"--image":    true,
		"--replicas": true,
		"--port":     true,
		"--timeout":  true,
		"--grace-period":    true,
		"--service-account": true,
		"--cluster":         true,
		"--user":            true,
		"--server":          true,
	}

	// Flags whose values collapse to <path>.
	pathFlags := map[string]bool{
		"-f": true, "--filename": true,
		"--kubeconfig": true,
		"--certificate-authority": true,
		"--client-certificate":    true,
		"--client-key":            true,
		"--cache-dir":             true,
	}

	// Boolean flags (no argument consumed).
	booleanFlags := map[string]bool{
		"-A": true, "--all-namespaces": true,
		"-w": true, "--watch": true,
		"--watch-only": true,
		"--all":        true,
		"--force":      true,
		"--recursive":  true,
		"--dry-run":    true,
		"--no-headers": true,
		"--show-labels": true,
		"-it": true, "-ti": true,
		"-i": true, "--stdin": true,
		"-t": true, "--tty": true,
		"--cascade":    true,
		"--overwrite":  true,
		"--privileged": true,
		"-d": true,
		"--prune":   true,
		"--record":  true,
		"--verbose": true,
		"-v":        true,
		"--version": true,
		"--help":    true,
	}

	var result []string
	resourceTypeAssigned := nameFirstSubcmds[subcommand] // skip resource type for name-first subcommands
	doubleDash := false

	i := 0
	for i < len(args) {
		tok := args[i]

		// After --, everything is passed through to the container command.
		// Keep those tokens verbatim (they are structural, e.g. "bash", "/bin/sh").
		if doubleDash {
			result = append(result, tok)
			i++
			continue
		}

		if tok == "--" {
			doubleDash = true
			result = append(result, tok)
			i++
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !resourceTypeAssigned {
				resourceTypeAssigned = true
			}
			i++
			continue
		}

		// Boolean flags: keep verbatim, no value consumed.
		if booleanFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Verbatim value flags: keep flag and value as-is.
		if verbatimValueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
			continue
		}

		// Selector flags: collapse value to <selector>.
		if selectorFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<selector>")
				}
				i++
			}
			continue
		}

		// Value flags: collapse value to <val>.
		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// Path flags: collapse value to <path>.
		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		// Other flags: keep verbatim.
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first one is the resource type, kept verbatim.
		if !resourceTypeAssigned {
			result = append(result, tok)
			resourceTypeAssigned = true
		} else {
			// Resource names are data; collapse to <word> unless classifyToken
			// gives a more specific placeholder (path, URL, etc.).
			classified := classifyToken(tok)
			if classified == tok {
				// classifyToken returned verbatim, so it's a plain word (resource name).
				result = append(result, "<word>")
			} else {
				result = append(result, classified)
			}
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
