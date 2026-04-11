package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("ansible", handleAnsibleAdhoc)
	shellshape.Register("ansible-playbook", handleAnsiblePlaybook)
}

var ansibleCategories = []shellshape.FlagCategory{
	{Flags: map[string]bool{
		"-a": true, "--args": true,
		"-e": true, "--extra-vars": true,
		"-u": true, "--user": true,
		"--become-user":     true,
		"--become-method":   true,
		"--tags":            true,
		"--skip-tags":       true,
		"--start-at":        true,
		"--start-at-task":   true,
		"--connection":      true,
		"--ssh-extra-args":  true,
		"--ssh-common-args": true,
		"--scp-extra-args":  true,
		"--sftp-extra-args": true,
	}, Placeholder: "<val>"},
	{Flags: map[string]bool{
		"-i": true, "--inventory": true, "--inventory-file": true,
		"--private-key": true, "--key-file": true,
		"--vault-password-file": true, "--vault-id": true,
	}, Placeholder: "<path>"},
	{Flags: map[string]bool{
		"-m": true, "--module-name": true,
	}, Placeholder: "<module>"},
	{Flags: map[string]bool{
		"-f": true, "--forks": true,
		"-t": true, "--timeout": true,
	}, Placeholder: "N"},
	{Flags: map[string]bool{
		"-l": true, "--limit": true,
	}, Placeholder: "<target>"},
}

// handleAnsibleAdhoc handles the `ansible` ad-hoc command.
// First positional is the host pattern (<target>), remaining positionals
// use classifyToken. Module flags (-m) collapse to <module>.
func handleAnsibleAdhoc(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	firstPositional := true

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			firstPositional = false
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, ansibleCategories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, ansibleCategories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional is the host pattern.
		if firstPositional {
			result = append(result, "<target>")
			firstPositional = false
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleAnsiblePlaybook handles the `ansible-playbook` command.
// All positionals are playbook files (classifyToken). Shares flag
// categories with the ad-hoc handler.
func handleAnsiblePlaybook(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, ansibleCategories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, ansibleCategories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// All positionals are playbook files.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
