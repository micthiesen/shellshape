package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("helm", handleHelm, shellshape.HandlerOptions{HasSubcommands: true})
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
	args, redirects := shellshape.SplitRedirects(tokens)

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

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: use classifyToken for generic classification
		// (URLs, numbers, paths get collapsed; bare words stay verbatim)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
