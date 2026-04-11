package shellshape

func init() {
	Register("zip", handleZip)
}

// handleZip handles the zip command.
// First positional is the archive name (<archive>).
// Remaining positionals are input files/paths (classifyToken).
// -b, --temp-path consume the next token as <path>.
// -P, --password consume the next token as <val>.
// -t, -tt consume the next token as <val> (date).
// -n consume the next token as <val> (suffixes).
// -Z, --compression-method consume the next token as <val>.
// -x, -i consume all following non-flag tokens as <pattern>.
func handleZip(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-b": true, "--temp-path": true,
	}

	valFlags := map[string]bool{
		"-P": true, "--password": true,
		"-t": true, "-tt": true,
		"-n": true,
		"-Z": true, "--compression-method": true,
	}

	// -x and -i start a list of patterns that continues until a flag is seen.
	listPatternFlags := map[string]bool{
		"-x": true, "--exclude": true,
		"-i": true, "--include": true,
	}

	// Boolean flags that isFlagToken won't catch (non-letter after dash).
	boolFlags := map[string]bool{
		"-0": true, "-1": true, "-2": true, "-3": true, "-4": true,
		"-5": true, "-6": true, "-7": true, "-8": true, "-9": true,
		"-@": true,
	}

	var result []string
	archiveAssigned := false
	inPatternList := false

	isZipFlag := func(tok string) bool {
		return isFlagToken(tok) || boolFlags[tok] || listPatternFlags[tok] ||
			pathFlags[tok] || valFlags[tok]
	}

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			inPatternList = false
			result = append(result, tok)
			if !archiveAssigned {
				archiveAssigned = true
			}
			i++
			continue
		}

		if listPatternFlags[tok] {
			inPatternList = true
			result = append(result, tok)
			i++
			continue
		}

		// While in a pattern list (-x/-i), non-flag tokens are patterns.
		if inPatternList {
			if isZipFlag(tok) {
				inPatternList = false
				// Fall through to flag handling below.
			} else {
				result = append(result, "<pattern>")
				i++
				continue
			}
		}

		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: first is archive, rest are paths.
		if !archiveAssigned {
			result = append(result, "<archive>")
			archiveAssigned = true
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
