package shellshape

import "strings"

func init() {
	Register("ruby", handleRuby)
}

// handleRuby handles the ruby interpreter.
// -e takes inline code. -r takes a library name. -I and -C take directory paths.
// -E/-F take values. Bundled short flags like -ne, -pe, -lne are split so that
// if -e is present, the next token is consumed as <code>.
// Fused flags like -Idir, -Cdir, -rlib, -F: are split into flag + value.
//
// Positional handling depends on context:
//   - With -n/-p (line loop) + -e: positionals are input file paths
//   - With -e alone (no line loop): positionals are script arguments (<arg>)
//   - Without -e: first positional is the script path, rest are <arg>
func handleRuby(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{"-I": true, "-C": true}
	valFlags := map[string]bool{
		"-r": true, "-E": true, "-F": true,
		"--encoding": true, "--external-encoding": true, "--internal-encoding": true,
	}
	codeFlags := map[string]bool{"-e": true}

	fusedArgFlags := map[byte]string{
		'I': "<path>", 'C': "<path>", 'r': "<val>", 'F': "<val>",
	}
	boolShortFlags := map[byte]bool{
		'a': true, 'c': true, 'd': true, 'l': true,
		'n': true, 'p': true, 's': true, 'S': true,
		'v': true, 'w': true,
	}

	var result []string
	hasCode := false     // -e was seen
	hasLineLoop := false // -n or -p was seen
	scriptSeen := false  // script path positional consumed (only for non -e mode)
	doubleDash := false

	// Pre-scan to detect -n/-p in bundled flags too.
	for _, tok := range args {
		if tok == "-n" || tok == "-p" {
			hasLineLoop = true
		}
		if len(tok) > 1 && tok[0] == '-' && tok[1] != '-' {
			for _, ch := range tok[1:] {
				if ch == 'n' || ch == 'p' {
					hasLineLoop = true
				}
			}
		}
	}

	i := 0
	for i < len(args) {
		tok := args[i]

		if tok == "--" {
			result = append(result, "--")
			doubleDash = true
			i++
			continue
		}

		// After --, everything is positional.
		if doubleDash {
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else if hasCode && hasLineLoop {
				result = append(result, classifyToken(tok))
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		// After script path in non-code mode, everything is an arg.
		if scriptSeen && !hasCode {
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !hasCode {
				scriptSeen = true
			}
			i++
			continue
		}

		// Exact match: code flags.
		if codeFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<code>")
			hasCode = true
			continue
		}

		// Exact match: path flags.
		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// Exact match: value flags.
		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Handle bundled short flags (-ne, -pe, -lne) and fused flags (-Ilib, -rjson).
		if len(tok) > 1 && tok[0] == '-' && tok[1] != '-' {
			parsed := parseRubyBundled(tok, boolShortFlags, fusedArgFlags)
			if parsed != "" {
				parts := strings.Fields(parsed)
				lastPart := parts[len(parts)-1]

				if lastPart == "<code>" {
					// Bundled flags include -e; next arg is the code.
					flagPart := strings.Join(parts[:len(parts)-1], " ")
					if flagPart != "" {
						result = append(result, flagPart)
					}
					i++
					if i < len(args) {
						if isSubshellToken(args[i]) {
							result = append(result, args[i])
						} else {
							result = append(result, "<code>")
						}
					}
					hasCode = true
					i++
					continue
				}

				result = append(result, parts...)
				i++
				continue
			}
		}

		// Long flags with = syntax.
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			eqIdx := strings.IndexByte(tok, '=')
			result = append(result, tok[:eqIdx]+"=<val>")
			i++
			continue
		}

		// Any other flag (boolean).
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional argument.
		if hasCode {
			// With -e code: positionals are file paths if line loop, else args.
			if hasLineLoop {
				result = append(result, classifyToken(tok))
			} else {
				result = append(result, "<arg>")
			}
		} else {
			// No -e: first positional is the script path.
			result = append(result, "<script>")
			scriptSeen = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// parseRubyBundled handles bundled short flags like -ne, -pe, -lne, -i.bak,
// and fused argument flags like -Ilib, -Cdir, -rjson, -F:.
// Returns a space-separated string of normalized tokens, or "" if not recognized.
// If the bundle contains -e, the last token will be "<code>" to signal that
// the caller should consume the next argument.
func parseRubyBundled(tok string, boolFlags map[byte]bool, fusedFlags map[byte]string) string {
	s := tok[1:] // strip leading -

	// Special case: -i with optional extension (e.g., -i.bak)
	if s[0] == 'i' {
		if len(s) == 1 {
			return "-i"
		}
		return tok
	}

	var parts []string
	j := 0
	for j < len(s) {
		ch := s[j]

		// Fused argument flag: -Ilib, -Cdir, etc.
		if placeholder, ok := fusedFlags[ch]; ok {
			flag := "-" + string(ch)
			if j+1 < len(s) {
				parts = append(parts, flag, placeholder)
				return strings.Join(parts, " ")
			}
			parts = append(parts, flag)
			return strings.Join(parts, " ")
		}

		// -e in bundle: emit bundled booleans + e, signal code.
		if ch == 'e' {
			bundled := "-"
			for _, p := range parts {
				bundled += p[1:]
			}
			bundled += "e"
			return bundled + " <code>"
		}

		if boolFlags[ch] {
			parts = append(parts, "-"+string(ch))
			j++
			continue
		}

		return ""
	}

	return strings.Join(parts, " ")
}
