package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("cargo", handleCargo, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleCargo handles the Rust cargo package manager.
// Package names and target triples are structural (kept verbatim).
// --features collapses its arg to <val>.
// -j/--jobs collapses its arg to N.
// --manifest-path, --target-dir collapse to <path>.
// Positionals are kept verbatim (package names, test filters).
func handleCargo(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Structural flags: their arguments stay verbatim (package names, targets, binaries)
	structuralFlags := map[string]bool{
		"-p": true, "--package": true,
		"--bin": true, "--example": true,
		"--test": true, "--bench": true,
		"--target": true,
	}

	pathFlags := map[string]bool{
		"--manifest-path": true,
		"--target-dir":    true,
	}

	numericFlags := map[string]bool{
		"-j": true, "--jobs": true,
	}

	valFlags := map[string]bool{
		"--features": true, "-F": true,
		"--color":   true,
		"--edition": true,
		"--profile": true,
		"--config":  true,
		"-Z":        true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		// Fused --flag=value for non-structural flags
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Structural flags: keep both flag and value verbatim
		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
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

		// Positional: keep verbatim (package names, test filters, etc.)
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
