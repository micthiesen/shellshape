package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"nvidia-smi", "nvidia-settings"} {
		shellshape.Register(name, handleNvidiaSmi)
	}
}

func handleNvidiaSmi(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-i": true, "-l": true,
	}

	// -d is structural: its argument is kept verbatim (MEMORY, TEMPERATURE, etc.)
	structuralFlags := map[string]bool{
		"-d": true,
	}

	// Fused flags: --query-gpu=<val> collapses, --format= and --id= are special
	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle --flag=value fused syntax
		if strings.Contains(tok, "=") && shellshape.IsFlagToken(tok) {
			eqIdx := strings.IndexByte(tok, '=')
			key := tok[:eqIdx]
			val := tok[eqIdx+1:]

			switch key {
			case "--query-gpu":
				result = append(result, key+"=<val>")
			case "--id":
				result = append(result, key+"=N")
			default:
				// Structural fused flags (--format=csv,noheader): keep verbatim
				_ = val
				result = append(result, tok)
			}
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if structuralFlags[tok] {
			// Keep next token verbatim
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional - classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
