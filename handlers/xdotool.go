package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"xdotool", "kdotool", "ydotool"} {
		shellshape.Register(name, handleXdotool, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

// handleXdotool handles xdotool, kdotool, and ydotool.
// The subcommand is already extracted by the normalizer.
//
// search: flags like --name, --class, --classname, --pid take a value arg
// that collapses to <val>. The trailing positional (pattern) also becomes <val>.
//
// type: typed text collapses to <str>. --delay and --window consume N.
//
// key/keyup/keydown: key names stay literal. --delay, --repeat, --window consume N.
//
// click/mousedown/mouseup: button number -> N. --repeat, --delay, --window consume N.
//
// mousemove/mousemove_relative: coordinates -> N. --window, --screen consume N.
//
// Window subcommands (windowfocus, windowactivate, etc.): window ID -> N,
// size/position args -> N.
func handleXdotool(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	switch subcommand {
	case "search":
		result = handleXdotoolSearch(args)
	case "type":
		result = handleXdotoolType(args)
	case "key", "keyup", "keydown":
		result = handleXdotoolKey(args)
	default:
		result = handleXdotoolDefault(args)
	}

	result = append(result, redirects...)
	return result
}

// handleXdotoolSearch normalizes xdotool search arguments.
// Value flags collapse their argument to <val>. The trailing positional
// (search pattern) also becomes <val>. Boolean flags stay literal.
func handleXdotoolSearch(args []string) []string {
	valFlags := map[string]bool{
		"--name": true, "--class": true, "--classname": true,
		"--pid": true, "--screen": true, "--desktop": true,
		"--limit": true, "--maxdepth": true,
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

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: the search pattern
		result = append(result, "<val>")
		i++
	}

	return result
}

// handleXdotoolType normalizes xdotool type arguments.
// --delay and --window consume N. All positional text becomes <str>.
func handleXdotoolType(args []string) []string {
	numericFlags := map[string]bool{
		"--delay": true, "--window": true,
	}

	var result []string
	hasText := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			if hasText {
				result = append(result, "<str>")
				hasText = false
			}
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional text: accumulate into single <str>
		hasText = true
		i++
	}

	if hasText {
		result = append(result, "<str>")
	}

	return result
}

// handleXdotoolKey normalizes xdotool key/keyup/keydown arguments.
// --delay, --repeat, --window consume N. Key names stay literal.
func handleXdotoolKey(args []string) []string {
	numericFlags := map[string]bool{
		"--delay": true, "--repeat": true, "--window": true,
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

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Key names stay literal
		result = append(result, tok)
		i++
	}

	return result
}

// handleXdotoolDefault handles all other xdotool subcommands generically.
// Numeric flags (--delay, --repeat, --window, --screen) consume N.
// Positionals get classifyToken (numbers become N, etc.).
func handleXdotoolDefault(args []string) []string {
	numericFlags := map[string]bool{
		"--delay": true, "--repeat": true, "--window": true, "--screen": true,
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

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	return result
}
