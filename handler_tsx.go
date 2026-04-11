package shellshape

func init() {
	Register("tsx", handleTsx)
}

// handleTsx handles the tsx (TypeScript Execute) command.
// Subcommands: watch.
// Flags --tsconfig consume a path argument.
// Flags -r/--require and --import consume a module argument.
// First positional is the script path → <script>.
// All positionals after the script are collapsed to <arg>.
func handleTsx(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{"--tsconfig": true}
	moduleFlags := map[string]bool{"-r": true, "--require": true, "--import": true}
	subcommands := map[string]bool{"watch": true}

	var result []string
	scriptSeen := false
	subcommandConsumed := false
	argCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		// Subcommand (only first positional, before script)
		if !subcommandConsumed && !scriptSeen && subcommands[tok] {
			result = append(result, tok)
			subcommandConsumed = true
			i++
			continue
		}

		if isSubshellToken(tok) {
			// Flush accumulated args before subshell
			for j := 0; j < argCount; j++ {
				result = append(result, "<arg>")
			}
			argCount = 0
			result = append(result, tok)
			scriptSeen = true
			i++
			continue
		}

		// Once we've seen the script, everything else is a script arg
		if scriptSeen {
			argCount++
			i++
			continue
		}

		// Flags that consume a path
		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		// Flags that consume a module specifier
		if moduleFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<module>")
				}
				i++
			}
			continue
		}

		// Boolean flags
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional = script
		result = append(result, "<script>")
		scriptSeen = true
		i++
	}

	// Flush remaining script args
	for j := 0; j < argCount; j++ {
		result = append(result, "<arg>")
	}

	result = append(result, redirects...)
	return result
}
