package shellshape

import "strings"

func init() {
	Register("screen", handleScreen)
}

// handleScreen handles the GNU screen terminal multiplexer.
//
// Value flags (-S, -s, -T, -e, -t, -p) collapse their argument to <val>.
// Path flags (-c) collapse their argument to <path>.
// Numeric flags (-h) collapse their argument to N.
// Boolean flags (-d, -D, -r, -R, -x, -X, -L, -ls, etc.) are kept verbatim.
// Compound flags like -dmS are recognized: if the compound ends with a
// value-consuming flag letter, the next token is consumed.
// Once the first positional is reached, all remaining tokens (the command to
// run inside screen, or session name) collapse to <val> or <path>.
func handleScreen(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-S": true, // session name
		"-s": true, // shell
		"-T": true, // terminal type
		"-e": true, // escape character
		"-t": true, // title
		"-p": true, // window number/name
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-c": true, // config file
	}

	// Numeric flags.
	numericFlags := map[string]bool{
		"-h": true, // scrollback lines
	}

	// Boolean flags (no argument consumed).
	boolFlags := map[string]bool{
		"-d": true, "-D": true, "-r": true, "-R": true,
		"-x": true, "-X": true, "-A": true, "-a": true,
		"-L": true, "-l": true, "-q": true, "-v": true,
		"-ls": true, "-list": true, "-wipe": true, "-m": true,
	}

	// Letters that consume the next token when at the end of a compound flag.
	valFlagLetters := map[byte]bool{
		'S': true, 's': true, 'T': true, 'e': true,
		't': true, 'p': true,
	}
	pathFlagLetters := map[byte]bool{
		'c': true,
	}
	numericFlagLetters := map[byte]bool{
		'h': true,
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

		// Exact match on known flags.
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

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Compound flags like -dmS, -dRR, etc.
		if strings.HasPrefix(tok, "-") && len(tok) > 2 && !strings.HasPrefix(tok, "--") {
			last := tok[len(tok)-1]
			if valFlagLetters[last] {
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
			if pathFlagLetters[last] {
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
			if numericFlagLetters[last] {
				result = append(result, tok)
				i++
				if i < len(args) {
					if isSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "N")
					}
					i++
				}
				continue
			}
			// Compound boolean flags like -dR, -DR
			result = append(result, tok)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional signals the start of either a session name or a
		// command to run inside screen. From here on, collapse everything
		// (including flags meant for the inner command) to <val> / <path>.
		for i < len(args) {
			tok = args[i]
			if isSubshellToken(tok) {
				result = append(result, tok)
			} else {
				classified := classifyToken(tok)
				if classified == tok {
					classified = "<val>"
				}
				result = append(result, classified)
			}
			i++
		}
	}

	result = append(result, redirects...)
	return result
}
