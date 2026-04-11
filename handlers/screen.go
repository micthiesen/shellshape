package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("screen", handleScreen)
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
	args, redirects := shellshape.SplitRedirects(tokens)

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

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Exact match on known flags.
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
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
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
				continue
			}
			if pathFlagLetters[last] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
				continue
			}
			if numericFlagLetters[last] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
				continue
			}
			// Compound boolean flags like -dR, -DR
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional signals the start of either a session name or a
		// command to run inside screen. From here on, collapse everything
		// (including flags meant for the inner command) to <val> / <path>.
		for i < len(args) {
			tok = args[i]
			if shellshape.IsSubshellToken(tok) {
				result = append(result, tok)
			} else {
				classified := shellshape.ClassifyToken(tok)
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
