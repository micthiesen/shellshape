package shellshape

func init() {
	Register("helm", handleHelm, HandlerOptions{HasSubcommands: true})
}

// handleHelm handles helm subcommand arguments.
// Since helm is in subcommandExecutables, the subcommand (install, upgrade,
// uninstall, list, repo, etc.) is already consumed before this handler is called.
//
// Positionals are kept verbatim (release names, chart references, repo names
// are all structural). Value flags (--set, --namespace, etc.) collapse their
// argument to <val>. Path flags (-f/--values, --kubeconfig) collapse to <path>.
// Boolean flags are kept verbatim.
func handleHelm(_ string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"--set": true, "--set-string": true, "--set-file": true,
		"--set-json": true,
		"-n":         true, "--namespace": true,
		"--version":      true,
		"--timeout":      true,
		"--kube-context": true,
		"--description":  true,
		"--output":       true, "-o": true,
		"--filter":      true,
		"--repo":        true,
		"--username":    true,
		"--password":    true,
		"--ca-file":     true,
		"--cert-file":   true,
		"--key-file":    true,
		"--history-max": true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-f": true, "--values": true,
		"--kubeconfig":        true,
		"--post-renderer":     true,
		"--registry-config":   true,
		"--repository-cache":  true,
		"--repository-config": true,
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
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

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

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: use classifyToken for generic classification
		// (URLs, numbers, paths get collapsed; bare words stay verbatim)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
