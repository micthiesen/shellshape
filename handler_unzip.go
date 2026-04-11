package shellshape

func init() {
	Register("unzip", handleUnzip)
}

// handleUnzip handles the unzip command.
// The first positional argument is the archive file (<archive>).
// Subsequent positionals are member patterns to extract (<pattern>).
// -d consumes the next token as <path> (destination directory).
// -P consumes the next token as <val> (password).
// -x switches to exclude mode: subsequent positionals become <pattern>.
// All other flags are boolean modifiers.
func handleUnzip(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

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

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !archiveAssigned {
				archiveAssigned = true
			}
			i++
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// -x: emit the flag then consume all following non-flag tokens as <pattern>.
		if tok == "-x" {
			result = append(result, "-x")
			i++
			for i < len(args) && !isFlagToken(args[i]) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<pattern>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
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
