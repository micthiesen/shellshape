package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"sensors", "sensors-detect"} {
		shellshape.Register(name, handleSensors)
	}
}

// handleSensors handles sensors and sensors-detect (lm_sensors).
// Boolean flags for output format. --bus takes a value, -c/--config-file takes a path.
// Positional chip names collapse to <val>.
func handleSensors(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{
			"--bus": true,
		}, Placeholder: "<val>"},
		{Flags: map[string]bool{
			"-c": true, "--config-file": true,
		}, Placeholder: "<path>"},
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

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: chip name → <val>
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
