package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("systemd-analyze", handleSystemdAnalyze, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleSystemdAnalyze handles systemd-analyze subcommand arguments.
// Since systemd-analyze has HasSubcommands, the normalizer extracts the first
// token after the executable as the "subcommand". However, global flags like
// --no-pager, --user, -H can appear before the subcommand, so this handler
// detects that case and finds the real subcommand from remaining tokens.
//
// Subcommands that take time expressions (calendar, timestamp, timespan) collapse
// positionals to <val>. The verify subcommand collapses positionals to <path>.
// Other subcommands treat positionals as unit names → <unit>.
func handleSystemdAnalyze(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-H": true, "--host": true,
		"-M": true, "--machine": true,
	}

	pathFlags := map[string]bool{
		"--root":  true,
		"--image": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}

	var result []string
	i := 0

	result, subcommand, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, categories, nil,
	)
	realSubcommand := subcommand

	// Determine placeholder for positionals based on subcommand.
	positionalPlaceholder := "<unit>"
	switch realSubcommand {
	case "calendar", "timestamp", "timespan":
		positionalPlaceholder = "<val>"
	case "verify":
		positionalPlaceholder = "<path>"
	}

	hasPositional := false

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

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional argument
		if !hasPositional {
			result = append(result, positionalPlaceholder)
			hasPositional = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
