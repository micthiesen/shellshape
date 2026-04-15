package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("strace", handleStrace)
	shellshape.Register("dtrace", handleDtrace)
}

// handleStrace handles strace (Linux system call tracer).
// Expression flags (-e, --trace, --signal, etc.) collapse to <expr>.
// Path flags (-o, -P, --output) collapse to <path>.
// Numeric flags (-p, -s, -a) collapse to N.
// Value flags (-S, -b, -I) collapse to <val>.
// All positionals (the traced command and its args) collapse to <cmd>.
func handleStrace(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-e": true,
		}, Placeholder: "<expr>"},
		{Flags: map[string]bool{
			"-o": true, "-P": true,
		}, Placeholder: "<path>"},
		{Flags: map[string]bool{
			"-p": true, "-s": true, "-a": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{
			"-S": true, "-b": true, "-I": true,
		}, Placeholder: "<val>"},
	}

	fusedCategories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"--trace": true, "--signal": true, "--status": true,
			"--abbrev": true, "--verbose": true, "--raw": true,
			"--read": true, "--write": true, "--fault": true,
			"--inject": true, "--kvm": true, "-e": true,
		}, Placeholder: "<expr>"},
		{Flags: map[string]bool{
			"--output": true,
		}, Placeholder: "<path>"},
		{Flags: map[string]bool{
			"--columns": true, "--sortby": true,
		}, Placeholder: "<val>"},
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

		if fused, ok := shellshape.ConsumeFusedFlag(tok, fusedCategories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional starts the traced command; collapse all remaining tokens.
		for i < len(args) {
			result = shellshape.EmitPositional(result, args[i], "<cmd>")
			i++
		}
		break
	}

	result = append(result, redirects...)
	return result
}

// handleDtrace handles dtrace (DTrace dynamic tracing).
// Probe flags (-n, -f, -i, -m) collapse to <probe>.
// Path flags (-s, -o) collapse to <path>.
// Numeric flag (-p) collapses to N.
// Command flag (-c) collapses to <cmd>.
// Value flags (-b, -D, -U, -x) collapse to <val>.
// Positionals collapse to <arg>.
func handleDtrace(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"-n": true, "-f": true, "-i": true, "-m": true,
		}, Placeholder: "<probe>"},
		{Flags: map[string]bool{
			"-s": true, "-o": true,
		}, Placeholder: "<path>"},
		{Flags: map[string]bool{
			"-p": true,
		}, Placeholder: "N"},
		{Flags: map[string]bool{
			"-c": true,
		}, Placeholder: "<cmd>"},
		{Flags: map[string]bool{
			"-b": true, "-D": true, "-U": true, "-x": true,
		}, Placeholder: "<val>"},
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
		result = append(result, "<arg>")
		i++
	}

	result = append(result, redirects...)
	return result
}
