package shellshape

func init() {
	Register("docker", handleDocker, HandlerOptions{HasSubcommands: true})
}

// handleDocker handles docker subcommand arguments.
// Since docker is in subcommandExecutables, the subcommand (run, build, exec,
// compose, etc.) is already consumed before this handler is called.
//
// Value flags (-v, -p, -e, --name, --network, etc.) collapse their argument to
// a placeholder. Path flags (-f, -w) collapse to <path>. Boolean flags are kept
// verbatim. Remaining positionals use classifyToken.
func handleDocker(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-v": true, "--volume": true,
		"-p": true, "--publish": true,
		"-e": true, "--env": true,
		"--name":    true,
		"--network": true, "--net": true,
		"-t": true, "--tag": true,
		"--build-arg":  true,
		"--entrypoint": true,
		"--platform":   true,
		"-u":           true, "--user": true,
		"-m": true, "--memory": true,
		"--cpus": true,
		"-l":     true, "--label": true,
		"--mount":      true,
		"--restart":    true,
		"--log-driver": true, "--log-opt": true,
		"--add-host": true,
		"--dns":      true, "--dns-search": true,
		"--expose":  true,
		"--link":    true,
		"--cap-add": true, "--cap-drop": true,
		"--device":       true,
		"--shm-size":     true,
		"--stop-signal":  true,
		"--stop-timeout": true,
		"--ulimit":       true,
		"--runtime":      true,
		"--pid":          true, "--ipc": true, "--uts": true,
		"--cgroupns": true,
		"--gpus":     true,
		"--pull":     true,
		"--ip":       true, "--ip6": true,
		"--hostname": true, "-h": true,
		"--mac-address":   true,
		"--domainname":    true,
		"--env-file":      true,
		"--network-alias": true,
		"--security-opt":  true,
		"--storage-opt":   true,
		"--sysctl":        true,
		"--tmpfs":         true,
		"--isolation":     true,
		"--format":        true,
		"--filter":        true,
		"--scale":         true,
		"--timeout":       true,
		"--target":        true,
		"--output":        true, "-o": true,
		"--progress":   true,
		"--cache-from": true, "--cache-to": true,
		"--secret": true, "--ssh": true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"-w": true, "--workdir": true,
	}

	categories := []flagCategory{
		{valFlags, "<val>"},
		{pathFlags, "<path>"},
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

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: classify generically
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
