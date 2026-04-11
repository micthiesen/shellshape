package shellshape

func init() {
	Register("strace", handleStrace)
	Register("dtrace", handleDtrace)
}

// handleStrace handles strace (Linux system call tracer).
// Expression flags (-e, --trace, --signal, etc.) collapse to <expr>.
// Path flags (-o, -P, --output) collapse to <path>.
// Numeric flags (-p, -s, -a) collapse to N.
// Value flags (-S, -b, -I) collapse to <val>.
// All positionals (the traced command and its args) collapse to <cmd>.
func handleStrace(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	categories := []flagCategory{
		{flags: map[string]bool{
			"-e": true,
		}, placeholder: "<expr>"},
		{flags: map[string]bool{
			"-o": true, "-P": true,
		}, placeholder: "<path>"},
		{flags: map[string]bool{
			"-p": true, "-s": true, "-a": true,
		}, placeholder: "N"},
		{flags: map[string]bool{
			"-S": true, "-b": true, "-I": true,
		}, placeholder: "<val>"},
	}

	fusedCategories := []flagCategory{
		{flags: map[string]bool{
			"--trace": true, "--signal": true, "--status": true,
			"--abbrev": true, "--verbose": true, "--raw": true,
			"--read": true, "--write": true, "--fault": true,
			"--inject": true, "--kvm": true, "-e": true,
		}, placeholder: "<expr>"},
		{flags: map[string]bool{
			"--output": true,
		}, placeholder: "<path>"},
		{flags: map[string]bool{
			"--columns": true, "--sortby": true,
		}, placeholder: "<val>"},
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

		if fused, ok := consumeFusedFlag(tok, fusedCategories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional starts the traced command; collapse all remaining tokens.
		for i < len(args) {
			if isSubshellToken(args[i]) {
				result = append(result, args[i])
			} else {
				result = append(result, "<cmd>")
			}
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
	args, redirects := splitRedirects(tokens)

	categories := []flagCategory{
		{flags: map[string]bool{
			"-n": true, "-f": true, "-i": true, "-m": true,
		}, placeholder: "<probe>"},
		{flags: map[string]bool{
			"-s": true, "-o": true,
		}, placeholder: "<path>"},
		{flags: map[string]bool{
			"-p": true,
		}, placeholder: "N"},
		{flags: map[string]bool{
			"-c": true,
		}, placeholder: "<cmd>"},
		{flags: map[string]bool{
			"-b": true, "-D": true, "-U": true, "-x": true,
		}, placeholder: "<val>"},
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

		// Positional
		result = append(result, "<arg>")
		i++
	}

	result = append(result, redirects...)
	return result
}
