package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("mkinitcpio", handleMkinitcpio)
}

// handleMkinitcpio handles mkinitcpio arguments.
// -p (preset) and -A/-S/-H (hook names) are structural (kept verbatim).
// -g/-c/-t (output/config/builddir paths) → <path>.
// -k (kernel version) → <val>.
func handleMkinitcpio(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-g": true, "--generate": true,
		"-c": true, "--config": true,
		"-t": true, "--builddir": true,
	}

	valFlags := map[string]bool{
		"-k": true, "--kernel": true,
	}

	// Structural flags: value is kept verbatim
	structuralFlags := map[string]bool{
		"-p": true, "--preset": true,
		"-A": true, "--addhooks": true,
		"-S": true, "--skiphooks": true,
		"-H": true, "--hookhelp": true,
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

		// Fused flags
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag categories with placeholder values
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
