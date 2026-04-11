package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("dmesg", handleDmesg)
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
	args, redirects := shellshape.SplitRedirects(tokens)

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

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: levelFlags, Placeholder: "<level>"},
		{Flags: facilityFlags, Placeholder: "<facility>"},
		{Flags: sizeFlags, Placeholder: "N"},
		{Flags: timeFlags, Placeholder: "<time>"},
		{Flags: fmtFlags, Placeholder: "<fmt>"},
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

		// Unexpected positional
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
