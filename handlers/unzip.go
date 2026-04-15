package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("unzip", handleUnzip)
}

// handleUnzip handles the unzip command.
// The first positional argument is the archive file (<archive>).
// Subsequent positionals are member patterns to extract (<pattern>).
// -d consumes the next token as <path> (destination directory).
// -P consumes the next token as <val> (password).
// -x switches to exclude mode: subsequent positionals become <pattern>.
// All other flags are boolean modifiers.
func handleUnzip(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"-d": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"-P": true,
	}

	var result []string
	archiveAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !archiveAssigned {
				archiveAssigned = true
			}
			i++
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

		// -x: emit the flag then consume all following non-flag tokens as <pattern>.
		if tok == "-x" {
			result = append(result, "-x")
			i++
			for i < len(args) && !shellshape.IsFlagToken(args[i]) {
				result = shellshape.EmitPositional(result, args[i], "<pattern>")
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first is archive, rest are member patterns.
		if !archiveAssigned {
			result = append(result, "<archive>")
			archiveAssigned = true
		} else {
			result = append(result, "<pattern>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
