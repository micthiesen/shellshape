package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("cpupower", handleCpupower, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleCpupower handles cpupower subcommand arguments.
// Governor names (-g) are structural. Frequency values (-f, -d, -u) → <val>.
// --cpu selector → <val>. Idle enable/disable values → <val>.
func handleCpupower(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose values collapse to <val>
	valFlags := map[string]bool{
		"-f": true, "--freq": true,
		"-d": true, "--min": true,
		"-u": true, "--max": true,
		"-e": true, "--enable": true,
		"-E": true, "--disable-by-latency": true,
		"-c": true, "--cpu": true,
	}

	// Structural flags: value kept verbatim (governor names)
	structuralFlags := map[string]bool{
		"-g": true, "--governor": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	// Handle --cpu or -c consumed as "subcommand" by the normalizer
	if shellshape.IsFlagToken(subcommand) {
		if valFlags[subcommand] {
			// Consume the value (e.g. "all", "0", "0-3")
			if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
				result = append(result, "<val>")
				i++
			}
		}
		// Emit the real subcommand
		if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) {
			result = append(result, args[i])
			i++
		}
	}

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused flags
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Value flags
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Structural flags: keep flag and value verbatim
		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Other flags: keep verbatim (boolean flags)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Unexpected positionals: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
