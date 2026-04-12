package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"which", "whereis", "type", "command"} {
		shellshape.Register(name, handleWhich)
	}
}

// handleWhich handles which, whereis, type, and command.
// All positional arguments (command names being looked up) collapse to <cmd>.
// whereis has -B, -M, -S flags that take a path argument.
func handleWhich(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// whereis flags that consume a path argument
	pathFlags := map[string]bool{"-B": true, "-M": true, "-S": true}

	var result []string
	hasCmdPositional := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			if hasCmdPositional {
				result = append(result, "<cmd>")
				hasCmdPositional = false
			}
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			if hasCmdPositional {
				result = append(result, "<cmd>")
				hasCmdPositional = false
			}
			if pathFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
				continue
			}
			result = append(result, tok)
			i++
			continue
		}

		// Positional: a command name being looked up
		hasCmdPositional = true
		i++
	}
	if hasCmdPositional {
		result = append(result, "<cmd>")
	}

	result = append(result, redirects...)
	return result
}
