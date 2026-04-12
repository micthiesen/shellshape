package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("mpv", handleMpv)
}

// handleMpv handles the mpv media player.
// mpv uses --key=value flags almost exclusively.
// --volume=N collapses the value to N.
// --start, --length collapse to <val>.
// --sub-file collapses to <path>.
// --profile keeps its value verbatim (structural choice).
// Boolean flags (--fullscreen, --no-video, etc.) have no value.
// Positionals are classified via ClassifyToken (paths or URLs).
func handleMpv(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"--volume": true}, Placeholder: "N"},
		{Flags: map[string]bool{"--start": true, "--length": true}, Placeholder: "<val>"},
		{Flags: map[string]bool{"--sub-file": true}, Placeholder: "<path>"},
	}

	// Flags whose fused value is kept verbatim
	verbatimFusedFlags := map[string]bool{
		"--profile": true,
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

		// Handle --flag=value syntax
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 && strings.HasPrefix(tok, "--") {
			key := tok[:eqIdx]

			// Verbatim fused flags: keep entire token as-is
			if verbatimFusedFlags[key] {
				result = append(result, tok)
				i++
				continue
			}

			// Category-matched fused flags
			if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, fused)
				i++
				continue
			}

			// Generic: ClassifyToken handles --flag=<val>
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: classify (path or URL)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
