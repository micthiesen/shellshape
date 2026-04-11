package shellshape

func init() {
	Register("deno", handleDeno, HandlerOptions{HasSubcommands: true})
}

// handleDeno handles the deno command.
// Deno is a subcommand-based tool (run, test, fmt, lint, compile, task, eval,
// serve, install, add, remove, check, info, bench, etc.). The subcommand is
// already extracted by normalizeSingleCommand.
//
// Path flags (--config, --import-map, --lock, --cert, --env-file, --output)
// consume the next token as <path>.
// Value flags (--target, --filter, --ext, --location, --inspect*) consume
// the next token as <val>.
// Numeric flags (--port) consume the next token as N.
// Permission flags (--allow-*, --deny-*) are boolean unless fused with =.
//
// Special subcommands:
//   - "task": first positional is the task name (structural), rest kept verbatim
//   - "eval": first positional is <code>, rest are <arg>
//   - "init": positionals kept verbatim (project name is structural)
func handleDeno(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// "task" subcommand: task name is structural, everything after is opaque
	if subcommand == "task" {
		var result []string
		for _, tok := range args {
			result = append(result, tok)
		}
		result = append(result, redirects...)
		return result
	}

	// "init" subcommand: project name is structural
	if subcommand == "init" {
		var result []string
		for _, tok := range args {
			result = append(result, tok)
		}
		result = append(result, redirects...)
		return result
	}

	pathFlags := map[string]bool{
		"-c": true, "--config": true,
		"--import-map": true, "--lock": true,
		"--cert": true, "--env-file": true,
		"--output": true, "-o": true,
		"--outdir": true, "--log-file": true,
		"--node-modules-dir": true,
	}

	valueFlags := map[string]bool{
		"--target": true, "--filter": true,
		"--ext": true, "--location": true,
		"--inspect": true, "--inspect-brk": true,
		"--inspect-wait": true,
		"--include":      true, "--exclude": true,
		"--v8-flags": true, "--seed": true,
		"--check": true,
	}

	numericFlags := map[string]bool{
		"--port": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{valueFlags, "<val>"},
		{numericFlags, "N"},
	}

	isEval := subcommand == "eval"

	var result []string
	codeSeen := false // for eval subcommand

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			if isEval && !codeSeen {
				codeSeen = true
			}
			i++
			continue
		}

		// Handle eval: first positional is code, rest are args
		if isEval && codeSeen {
			result = append(result, "<arg>")
			i++
			continue
		}

		// Fused flags (--flag=value)
		if fused, ok := consumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag categories
		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Other flags (boolean, or unknown --flag=value handled by classifyToken)
		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional
		if isEval {
			result = append(result, "<code>")
			codeSeen = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
