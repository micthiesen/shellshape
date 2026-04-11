package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("env", handleEnv)
}

// handleEnv handles the env command.
// Flags -u/--unset consume the next arg as <var>, -C/--chdir as <path>,
// -S/--split-string as <str>. KEY=VALUE assignments collapse to a single
// <assign>. The first non-assignment positional (the command to run) is kept
// verbatim. Remaining positionals use classifyToken.
func handleEnv(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-u": true, "--unset": true}, Placeholder: "<var>"},
		{Flags: map[string]bool{"-C": true, "--chdir": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"-S": true, "--split-string": true}, Placeholder: "<str>"},
	}

	var result []string
	hasAssign := false
	commandSeen := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			if hasAssign {
				result = append(result, "<assign>")
				hasAssign = false
			}
			result = append(result, tok)
			commandSeen = true
			i++
			continue
		}

		// After the command token, classify remaining args generically.
		if commandSeen {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Fused long Flags: --unset=HOME, --chdir=/tmp
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag with argument: -u VAR, -C /tmp, -S '...'
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Boolean flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// KEY=VALUE assignment
		if shellshape.IsEnvAssignment(tok) {
			hasAssign = true
			i++
			continue
		}

		// First non-flag, non-assignment positional is the command to run.
		if hasAssign {
			result = append(result, "<assign>")
			hasAssign = false
		}
		result = append(result, tok) // keep command name verbatim
		commandSeen = true
		i++
	}

	// Trailing assignments with no command (e.g. `env FOO=bar`)
	if hasAssign {
		result = append(result, "<assign>")
	}

	result = append(result, redirects...)
	return result
}
