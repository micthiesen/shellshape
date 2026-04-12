package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("claude", handleClaude)
}

// handleClaude handles the Claude Code CLI.
// Subcommands (mcp, config) are detected manually since flags can precede them.
// --model keeps its argument verbatim (structural).
// --allowedTools, --resume, --output-format collapse to <val>.
// --max-turns collapses to N.
// Positional prompt text collapses to <str>.
// For mcp/config subcommands, remaining positionals are kept verbatim.
func handleClaude(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	subcommands := map[string]bool{
		"mcp": true, "config": true, "chat": true,
	}

	structuralFlags := map[string]bool{
		"--model": true,
	}

	valFlags := map[string]bool{
		"--allowedTools": true, "--disallowedTools": true,
		"--resume": true, "--output-format": true,
		"--append-system-prompt": true, "--system-prompt": true,
		"--permission-mode": true,
	}

	numericFlags := map[string]bool{
		"--max-turns": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	inSubcommand := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Detect subcommand as first non-flag token.
		if !inSubcommand && subcommands[tok] {
			result = append(result, tok)
			inSubcommand = true
			i++
			// After a subcommand, keep remaining positionals verbatim
			for i < len(args) {
				result = append(result, args[i])
				i++
			}
			break
		}

		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
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

		// Positional: prompt text collapses to <str>
		result = append(result, "<str>")
		i++
	}

	result = append(result, redirects...)
	return result
}
