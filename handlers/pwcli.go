package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("pw-cli", handlePwCli, shellshape.HandlerOptions{HasSubcommands: true})
}

// handlePwCli handles pw-cli subcommand arguments.
// Subcommands: info, list, connect, disconnect, set-param, enum-params,
// dump, create-node, destroy-node.
//
// For list/dump: positionals are structural type names.
// For set-param: first positional (ID) → <val>, second (param name) → structural, rest → <val>.
// For enum-params: first positional (ID) → <val>, second (param name) → structural.
// For create-node: first positional (factory) → structural, rest → <val>.
// For everything else: positionals → <val>.
func handlePwCli(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	posIdx := 0
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			posIdx++
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional handling depends on subcommand
		switch subcommand {
		case "list", "dump":
			// Type names are structural
			result = append(result, tok)
		case "set-param":
			// pos 0: ID → <val>, pos 1: param name → structural, pos 2+: value → <val>
			if posIdx == 1 {
				result = append(result, tok)
			} else {
				result = append(result, "<val>")
			}
		case "enum-params":
			// pos 0: ID → <val>, pos 1: param name → structural
			if posIdx == 1 {
				result = append(result, tok)
			} else {
				result = append(result, "<val>")
			}
		case "create-node":
			// pos 0: factory name → structural, pos 1+: props → <val>
			if posIdx == 0 {
				result = append(result, tok)
			} else {
				result = append(result, "<val>")
			}
		default:
			result = append(result, "<val>")
		}
		posIdx++
		i++
	}

	result = append(result, redirects...)
	return result
}
