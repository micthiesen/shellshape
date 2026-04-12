package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("gdbus", handleGdbus, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleGdbus handles gdbus subcommand arguments.
// Subcommands: call, emit, introspect, monitor, wait.
// --dest and bus names collapse to <val>, --object-path to <path>,
// --method/--signal use ClassifyToken (dotted names become <dotted-id>).
// Remaining positionals (signal/method arguments) become <val>.
func handleGdbus(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"--dest": true, "-d": true,
	}
	pathFlags := map[string]bool{
		"--object-path": true, "-o": true,
		"--xml": true,
	}
	classifyFlags := map[string]bool{
		"--method": true, "--signal": true,
	}
	numericFlags := map[string]bool{
		"--timeout": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		// --method/--signal: consume next arg with ClassifyToken
		if classifyFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, shellshape.ClassifyToken(args[i]))
				}
				i++
			}
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: method/signal arguments → <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
