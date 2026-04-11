package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"fly", "flyctl"} {
		shellshape.Register(name, handleFly, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

// handleFly handles fly/flyctl (Fly.io CLI) commands.
// fly is in subcommandExecutables, so the first subcommand (deploy, launch,
// apps, machine, secrets, etc.) is already consumed before this handler.
//
// Many fly commands have a second subcommand (apps list, secrets set, machine
// stop, scale count) which is kept verbatim as structural.
//
// Value flags (--app, --region, --image, --env, etc.) collapse to <val>.
// Path flags (--dockerfile, --config, --into) collapse to <path>.
// Boolean flags (--detach, --no-cache, --yes, etc.) are kept verbatim.
// Remaining positionals after the second subcommand use classifyToken.
func handleFly(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-a": true, "--app": true,
		"-o": true, "--org": true,
		"-r": true, "--region": true,
		"-i": true, "--image": true,
		"-e": true, "--env": true,
		"-s": true, "--signal": true,
		"-t": true, "--access-token": true,
		"-v": true, "--volume": true,
		"--name":                    true,
		"--strategy":                true,
		"--build-arg":               true,
		"--build-secret":            true,
		"--build-target":            true,
		"--buildpacks-docker-host":  true,
		"--buildpacks-volume":       true,
		"--compression":             true,
		"--deploy-retries":          true,
		"--depot":                   true,
		"--depot-scope":             true,
		"--exclude-machines":        true,
		"--exclude-regions":         true,
		"--file-literal":            true,
		"--file-local":              true,
		"--file-secret":             true,
		"--image-label":             true,
		"--label":                   true,
		"--lease-timeout":           true,
		"--max-concurrent":          true,
		"--max-unavailable":         true,
		"--only-machines":           true,
		"--primary-region":          true,
		"--process-groups":          true,
		"--regions":                 true,
		"--release-command-timeout": true,
		"--vm-cpu-kind":             true,
		"--vm-cpus":                 true,
		"--vm-gpu-kind":             true,
		"--vm-gpus":                 true,
		"--vm-memory":               true,
		"--vm-size":                 true,
		"--volume-initial-size":     true,
		"--wait-timeout":            true,
		"--auto-stop":               true,
		"--command":                 true,
		"--db":                      true,
		"--from":                    true,
		"--host-dedication-id":      true,
		"--internal-port":           true,
		"--secret":                  true,
		"--format":                  true,
		"--select":                  true,
		"--size":                    true,
	}

	pathFlags := map[string]bool{
		"-c": true, "--config": true,
		"--dockerfile": true,
		"--path":       true,
		"--into":       true,
		"--ignorefile": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	// Subcommands where the first positional is NOT a second subcommand.
	// For these, all positionals are data and should be classified.
	noSecondSubcmd := map[string]bool{
		"deploy":    true,
		"launch":    true,
		"status":    true,
		"logs":      true,
		"dashboard": true,
		"open":      true,
		"version":   true,
		"doctor":    true,
		"console":   true,
		"proxy":     true,
		"info":      true,
		"ping":      true,
		"docs":      true,
	}

	var result []string
	subcommandTaken := noSecondSubcmd[subcommand]

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
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

		// Positional: first positional is the second subcommand, kept verbatim
		if !subcommandTaken {
			result = append(result, tok)
			subcommandTaken = true
			i++
			continue
		}

		// Remaining positionals: classify
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
