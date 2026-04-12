package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("coredumpctl", handleCoredumpctl, shellshape.HandlerOptions{HasSubcommands: true})
}

var pidRE = regexp.MustCompile(`^\d+$`)
var numericFlagRE = regexp.MustCompile(`^-\d+$`)

// handleCoredumpctl handles coredumpctl subcommand arguments.
// PIDs (numeric positionals) become N, executable names/patterns become <val>.
// --since/--until timestamps become <val>, --output/-o become <path>.
func handleCoredumpctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
		"-D": true, "--directory": true,
	}

	valFlags := map[string]bool{
		"-S": true, "--since": true,
		"-U": true, "--until": true,
		"--debugger": true,
		"-F":         true, "--field": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
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

		// Fused flags (--output=/tmp/core)
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag categories with arguments
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Numeric flags like -1 (IsFlagToken doesnt match these)
		if numericFlagRE.MatchString(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Other flags: keep verbatim (boolean flags like -r, -q, -1, --no-pager, etc.)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: PIDs become N, everything else becomes <val>
		if pidRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
