package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("makepkg", handleMakepkg)
}

// handleMakepkg handles makepkg (Arch Linux package build tool).
//
// Most flags are boolean (-s, -i, -c, -f, -r, -C, -L, -d, --noconfirm, etc.).
// -p/--pkgbuild and --config take a path argument.
// --key takes a generic value argument.
// Positional arguments are rare and treated as paths.
func handleMakepkg(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-p":       true,
		"--config": true,
	}

	valFlags := map[string]bool{
		"--key": true,
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional arguments (rare) - classify as paths
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
