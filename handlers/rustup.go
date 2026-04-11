package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("rustup", handleRustup, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleRustup handles the rustup (Rust toolchain installer) command.
// Since HasSubcommands is set, the first subcommand (install, update, default,
// toolchain, target, component, override, run, which, doc, self, etc.) is
// already extracted.
//
// Subcommands with a second subcommand (toolchain, target, component, override,
// self, show): the second token is kept verbatim as structural.
//
// Toolchain names (stable, nightly, 1.70.0) collapse to <toolchain>.
// Target triples collapse to <target>.
// Component names collapse to <component>.
// --toolchain flag consumes next token as <toolchain>.
// --path flag consumes next token as <path>.
// --profile flag consumes next token as <val>.
// For "run": first positional is <toolchain>, rest stay verbatim (command to execute).
// For "toolchain link": first positional is <toolchain>, second is <path>.
func handleRustup(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"--profile": true,
	}

	toolchainFlags := map[string]bool{
		"--toolchain": true,
	}

	pathFlags := map[string]bool{
		"--path": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: toolchainFlags, Placeholder: "<toolchain>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	// Subcommands that have a second subcommand.
	hasSecondSubcmd := map[string]bool{
		"toolchain": true,
		"target":    true,
		"component": true,
		"override":  true,
		"self":      true,
		"show":      true,
	}

	// Subcommands where positionals are toolchain names.
	toolchainPositional := map[string]bool{
		"install":   true,
		"uninstall": true,
		"update":    true,
		"default":   true,
	}

	var result []string
	secondSubcmdTaken := !hasSecondSubcmd[subcommand]
	var secondSubcmd string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
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

		// Positional argument handling.
		positionalCount++

		// Second subcommand: keep verbatim.
		if !secondSubcmdTaken {
			secondSubcmd = tok
			result = append(result, tok)
			secondSubcmdTaken = true
			i++
			continue
		}

		// "run" subcommand: first positional is toolchain, rest are the command (verbatim).
		if subcommand == "run" {
			if positionalCount == 1 {
				result = append(result, "<toolchain>")
			} else {
				// Remaining args are the command to run, keep verbatim.
				result = append(result, tok)
			}
			i++
			continue
		}

		// "which" subcommand: positional is a component name.
		if subcommand == "which" {
			result = append(result, "<component>")
			i++
			continue
		}

		// Subcommands with second subcommand.
		if subcommand == "toolchain" {
			switch secondSubcmd {
			case "install", "remove", "uninstall":
				result = append(result, "<toolchain>")
			case "link":
				if positionalCount == 2 {
					result = append(result, "<toolchain>")
				} else {
					result = append(result, shellshape.ClassifyToken(tok))
				}
			default:
				result = append(result, shellshape.ClassifyToken(tok))
			}
			i++
			continue
		}

		if subcommand == "target" {
			switch secondSubcmd {
			case "add", "remove":
				result = append(result, "<target>")
			default:
				result = append(result, shellshape.ClassifyToken(tok))
			}
			i++
			continue
		}

		if subcommand == "component" {
			switch secondSubcmd {
			case "add", "remove":
				result = append(result, "<component>")
			default:
				result = append(result, shellshape.ClassifyToken(tok))
			}
			i++
			continue
		}

		if subcommand == "override" {
			switch secondSubcmd {
			case "set":
				result = append(result, "<toolchain>")
			default:
				result = append(result, shellshape.ClassifyToken(tok))
			}
			i++
			continue
		}

		// Toolchain-positional subcommands (install, update, default, uninstall).
		if toolchainPositional[subcommand] {
			result = append(result, "<toolchain>")
			i++
			continue
		}

		// Everything else: keep verbatim (completions args, self subcommands, etc.)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
