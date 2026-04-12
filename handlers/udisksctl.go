package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("udisksctl", handleUdisksctl, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleUdisksctl handles udisksctl disk management commands.
// Subcommands (mount, unmount, info, etc.) are consumed by the framework.
// Path flags: -b/--block-device, -p/--object-path, -f/--file → <path>
// Value flags: -t/--filesystem-type, -o/--options → <val>
// Boolean flags: --no-user-interaction, --read-only, --force, etc.
func handleUdisksctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-b": true, "--block-device": true,
		"-p": true, "--object-path": true,
		"-f": true, "--file": true,
	}

	valFlags := map[string]bool{
		"-t": true, "--filesystem-type": true,
		"-o": true, "--options": true,
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

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
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

		// Positional: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
