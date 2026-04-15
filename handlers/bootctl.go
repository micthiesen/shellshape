package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("bootctl", handleBootctl, shellshape.HandlerOptions{HasSubcommands: true})
}

var bootctlNumericRE = regexp.MustCompile(`^\d+$`)

// handleBootctl handles bootctl subcommand arguments.
// Boot entry IDs → <val>, timeout numeric values → N, --esp-path/--boot-path → <path>.
func handleBootctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"--esp-path":  true,
		"--boot-path": true,
		"-p":          true,
	}

	valFlags := map[string]bool{
		"--entry-token": true,
		"--image":       true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	// Handle a leading flag that the normalizer extracted as the subcommand.
	result, _, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, categories, nil,
	)

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused flags (--esp-path=/efi)
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

		// Other flags: keep verbatim
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: numeric → N, otherwise → <val>
		if bootctlNumericRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
