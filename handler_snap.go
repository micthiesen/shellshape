package shellshape

func init() {
	Register("snap", handleSnap, HandlerOptions{HasSubcommands: true})
}

// handleSnap handles snap subcommand arguments.
// Since snap is in subcommandExecutables, the subcommand (install, remove,
// refresh, find, etc.) is already consumed before this handler is called.
//
// For install/remove/refresh/revert/download/enable/disable/info/list:
// positionals are snap names, kept verbatim (structural).
// For find: the positional is a search query, collapsed to <query>.
// For set: first positional is snap name (verbatim), rest are key=value pairs (<val>).
// For get: all positionals are snap name + keys, kept verbatim.
// For connect/disconnect: positionals are plug/slot references, kept verbatim.
// For change/watch/abort: positional is a numeric ID, collapsed to N.
// For ack: positional is a file path, classified generically.
// Flags: --channel and --revision consume a value; others are boolean.
func handleSnap(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Subcommands where positionals are snap package names (structural).
	packageSubcommands := map[string]bool{
		"install": true, "remove": true, "refresh": true, "revert": true,
		"download": true, "enable": true, "disable": true,
		"info": true, "list": true,
	}

	// Subcommands where positionals are interface references (structural).
	interfaceSubcommands := map[string]bool{
		"connect": true, "disconnect": true,
	}

	// Subcommands where the positional is a numeric change ID.
	numericIDSubcommands := map[string]bool{
		"change": true, "watch": true, "abort": true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"--channel": true,
	}

	// Flags that consume the next token as a numeric value.
	numFlags := map[string]bool{
		"--revision": true,
	}

	var result []string
	i := 0
	// For the "set" subcommand, track whether we've seen the snap name yet.
	seenSetName := false

	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle --flag=value forms.
		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if numFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional handling depends on subcommand.
		if subcommand == "find" {
			result = append(result, "<query>")
			i++
			continue
		}

		if subcommand == "set" {
			if !seenSetName {
				// First positional is the snap name (structural).
				result = append(result, tok)
				seenSetName = true
			} else {
				// Remaining positionals are key=value settings.
				result = append(result, "<val>")
			}
			i++
			continue
		}

		if subcommand == "ack" {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		if numericIDSubcommands[subcommand] {
			result = append(result, "N")
			i++
			continue
		}

		if packageSubcommands[subcommand] || interfaceSubcommands[subcommand] || subcommand == "get" {
			// Package names, interface references, and config keys are structural.
			result = append(result, tok)
			i++
			continue
		}

		// Unknown subcommand: generic classification.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
