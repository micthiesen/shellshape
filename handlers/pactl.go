package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("pactl", handlePactl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handlePactl handles pactl subcommand arguments.
// Subcommands: list, info, stat, set-sink-volume, set-source-volume,
// set-sink-mute, set-source-mute, move-sink-input, move-source-output,
// load-module, unload-module, set-default-sink, set-default-source.
//
// For list: positionals are structural (sink, source, module type names).
// For load-module: first positional (module name) is structural, rest are <val>.
// For all others: positionals are data → <val>.
// -s and -n flags consume next arg as <val>.
func handlePactl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-s": true, "--server": true,
		"-n": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	moduleNameSeen := false
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

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional handling depends on subcommand
		switch subcommand {
		case "list":
			// list args are structural (type names like sinks, sources, short)
			result = append(result, tok)
		case "load-module":
			if !moduleNameSeen {
				// First positional is the module name: structural
				result = append(result, tok)
				moduleNameSeen = true
			} else {
				result = append(result, "<val>")
			}
		default:
			// All other subcommands: positionals are data
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
