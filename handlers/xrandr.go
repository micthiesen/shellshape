package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("xrandr", handleXrandr)
}

// handleXrandr handles xrandr display configuration.
// --output, --mode, --pos, --left-of, --right-of, --above, --below, --same-as → <val>
// --rotate keeps its argument verbatim (structural: normal/left/right/inverted)
// --rate, --screen, --dpi → N
// --query, --verbose, --dryrun, --off, --auto, --primary, etc. are boolean.
func handleXrandr(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"--output": true, "--mode": true, "--pos": true,
			"--left-of": true, "--right-of": true, "--above": true,
			"--below": true, "--same-as": true, "--set": true,
			"--addmode": true, "--delmode": true, "--newmode": true,
			"--rmmode": true, "--fb": true, "--filter": true,
			"--transform": true, "--scale": true, "--scale-from": true,
			"--panning": true, "--gamma": true, "--brightness": true,
			"--reflect": true, "--crtc": true,
		}, Placeholder: "<val>"},
		{Flags: map[string]bool{
			"--rate": true, "--screen": true, "--dpi": true,
		}, Placeholder: "N"},
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

		// --rotate keeps its argument verbatim (structural directions)
		if tok == "--rotate" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
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

		// Positional args are unusual for xrandr; classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
