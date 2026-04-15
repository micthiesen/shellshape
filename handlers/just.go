package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("just", handleJust)
}

// handleJust handles the just command runner.
// Path flags (-f, -d, -E, --justfile, --working-directory, etc.) collapse their
// argument to <path>. Value flags (--shell, --color, --set, --completions, etc.)
// collapse to <val>. --set consumes TWO tokens (variable name and value).
// Boolean flags (--dry-run, --verbose, -l, --init, etc.) are kept verbatim.
// The first positional becomes <recipe>, remaining positionals become <arg>.
func handleJust(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-f": true, "--justfile": true,
		"-d": true, "--working-directory": true,
		"-E": true, "--dotenv-path": true,
		"--ceiling": true,
		"--cygpath": true,
		"--tempdir": true,
	}

	valFlags := map[string]bool{
		"--shell": true, "--shell-arg": true,
		"--color": true, "--command-color": true,
		"--alias-style":      true,
		"--chooser":          true,
		"--dotenv-filename":  true,
		"--dump-format":      true,
		"--evaluate-format":  true,
		"--justfile-name":    true,
		"--list-heading":     true,
		"--list-prefix":      true,
		"--group":            true,
		"--indentation":      true,
		"--timestamp-format": true,
		"--completions":      true,
		"-c":                 true, "--command": true,
		"-s": true, "--show": true,
		"--usage": true,
	}

	// --set takes TWO arguments: variable name and value.
	setFlags := map[string]bool{
		"--set": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: setFlags, Placeholder: "<val>"},
	}

	var result []string
	recipeAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			recipeAssigned = true
			i++
			continue
		}

		// Fused --flag=value
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			// --set=VAR still needs to consume the value argument
			eqIdx := len(tok)
			for j := 0; j < len(tok); j++ {
				if tok[j] == '=' {
					eqIdx = j
					break
				}
			}
			key := tok[:eqIdx]
			if setFlags[key] {
				i++
				if i < len(args) {
					result = shellshape.EmitPositional(result, args[i], "<val>")
					i++
				}
			} else {
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if setFlags[tok] {
			// --set takes two arguments: variable name and value
			result = append(result, tok)
			i++
			for j := 0; j < 2 && i < len(args); j++ {
				result = shellshape.EmitPositional(result, args[i], "<val>")
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !recipeAssigned {
			result = append(result, "<recipe>")
			recipeAssigned = true
		} else {
			result = append(result, "<arg>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
