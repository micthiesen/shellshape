package shellshape

func init() {
	Register("dmesg", handleDmesg)
}

// handleDmesg handles the dmesg command (kernel ring buffer).
// Flags that consume arguments:
//
//	-f/--facility    → <facility>
//	-l/--level       → <level>
//	-n/--console-level → <level>
//	-F/--file, -K/--kmsg-file, -M, -N → <path>
//	-s/--buffer-size → N
//	--since/--until  → <time>
//	--time-format    → <fmt>
//
// All other flags are boolean. No positional arguments expected;
// any stray positionals get classifyToken.
func handleDmesg(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-F": true, "--file": true,
		"-K": true, "--kmsg-file": true,
		"-M": true, "-N": true,
	}
	levelFlags := map[string]bool{
		"-l": true, "--level": true,
		"-n": true, "--console-level": true,
	}
	facilityFlags := map[string]bool{
		"-f": true, "--facility": true,
	}
	sizeFlags := map[string]bool{
		"-s": true, "--buffer-size": true,
	}
	timeFlags := map[string]bool{
		"--since": true, "--until": true,
	}
	fmtFlags := map[string]bool{
		"--time-format": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{levelFlags, "<level>"},
		{facilityFlags, "<facility>"},
		{sizeFlags, "N"},
		{timeFlags, "<time>"},
		{fmtFlags, "<fmt>"},
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

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Unexpected positional
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
