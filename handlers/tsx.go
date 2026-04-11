package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("tsx", handleTsx)
}

// handleTsx handles the tsx (TypeScript Execute) command.
// Subcommands: watch.
// Flags --tsconfig consume a path argument.
// Flags -r/--require and --import consume a module argument.
// First positional is the script path → <script>.
// All positionals after the script are collapsed to <arg>.
func handleTsx(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

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

		if shellshape.IsSubshellToken(tok) {
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
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// Flags that consume a module specifier
		if moduleFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<module>")
			continue
		}

		// Boolean flags
		if shellshape.IsFlagToken(tok) {
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
