package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("tmux", handleTmux, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleTmux handles tmux subcommand arguments.
// Since tmux is in subcommandExecutables, the subcommand (new-session, attach,
// send-keys, split-window, etc.) is already consumed before this handler runs.
//
// Value flags (-s, -t, -n, -F, -L, -b, -E) collapse their argument to <val>.
// Path flags (-c) collapse their argument to <path>.
// Boolean flags (-d, -h, -v, -r, -D, -P) are kept verbatim.
// For send-keys: remaining positionals collapse to <str> (the keys to send).
// For other subcommands: positionals use classifyToken.
func handleTmux(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-s": true, // session name
		"-t": true, // target (session, window, or pane)
		"-n": true, // window name
		"-F": true, // format string
		"-L": true, // socket name
		"-b": true, // buffer name
		"-E": true, // environment variable
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-c": true, // start directory
	}

	// Numeric flags.
	numericFlags := map[string]bool{
		"-x": true, // width
		"-y": true, // height
		"-l": true, // size
	}

	sendKeysMode := subcommand == "send-keys"

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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if sendKeysMode {
			result = append(result, "<str>")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
