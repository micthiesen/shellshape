package shellshape

import "strings"

func init() {
	Register("rpm", handleRpm)
}

// handleRpm handles rpm (RPM Package Manager) commands.
//
// rpm uses flag-based modes rather than subcommands:
//   - Install/Upgrade (-i, -U, --install, --upgrade): positionals are .rpm files -> <path>
//   - Checksig (-K, --checksig): positionals are .rpm files -> <path>
//   - Query file (-qf, --file): positionals are filesystem paths -> <path>
//   - Query package file (-qp, --package): positionals are .rpm files -> <path>
//   - Build (-ba, -bb, etc.): positionals are spec files -> <path>
//   - Erase (-e, --erase): positionals are package names -> verbatim
//   - Query (-q, --query): positionals are package names -> verbatim
//   - Verify (-V, --verify): positionals are package names -> verbatim
//
// Flags like --queryformat, --qf, --root, --dbpath, --define consume the next
// token as a value or path.
func handleRpm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"--root":   true,
		"--dbpath": true,
		"--rcfile": true,
		"--macros": true,
		"--prefix": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"--queryformat": true,
		"--qf":          true,
		"--define":      true,
		"--relocate":    true,
		"--pipe":        true,
	}

	// Long flags that indicate positionals are file paths.
	longFileMode := map[string]bool{
		"--install":  true,
		"--upgrade":  true,
		"--checksig": true,
		"--package":  true,
		"--file":     true,
	}

	// Determine the mode from all flags first (two-pass approach).
	// In file mode, positionals are paths; otherwise they're package names (verbatim).
	fileMode := false
	for _, tok := range args {
		if longFileMode[tok] {
			fileMode = true
			break
		}
		// Check bundled short flags for mode-indicating characters.
		// -i (install), -U (upgrade), -K (checksig) mean file mode.
		// -f (file query), -p (package file query) in combination with -q mean file mode.
		if len(tok) > 1 && tok[0] == '-' && tok[1] != '-' {
			chars := tok[1:]
			if strings.ContainsAny(chars, "iUK") {
				fileMode = true
				break
			}
			if strings.ContainsAny(chars, "fp") && strings.Contains(chars, "q") {
				fileMode = true
				break
			}
			// Also handle -f or -p without -q bundled, when -q appears as
			// a separate flag elsewhere. We'll check for standalone -q below.
		}
	}

	// Second pass: if we saw standalone -f or -p, check for standalone -q.
	if !fileMode {
		hasQ := false
		hasFP := false
		for _, tok := range args {
			if len(tok) > 1 && tok[0] == '-' && tok[1] != '-' {
				chars := tok[1:]
				if strings.Contains(chars, "q") {
					hasQ = true
				}
				if strings.ContainsAny(chars, "fp") {
					hasFP = true
				}
			}
		}
		if hasQ && hasFP {
			fileMode = true
		}
	}

	// Also detect build mode: -ba, -bb, -bp, -bi, -bl, -bs flags.
	if !fileMode {
		for _, tok := range args {
			if len(tok) == 3 && tok[0] == '-' && tok[1] == 'b' {
				fileMode = true
				break
			}
		}
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
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

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional argument.
		if fileMode {
			result = append(result, "<path>")
		} else {
			// Package names are structural, kept verbatim.
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
