package handlers

import (
	"regexp"

	shellshape "github.com/micthiesen/shellshape"
)

var durationRE = regexp.MustCompile(`^-?\d+(\.\d+)?(s|m|h|d)?$`)

func init() {
	shellshape.Register("timeout", handleTimeout)
}

func handleTimeout(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	signalFlags := map[string]bool{
		"-s": true, "--signal": true,
	}
	numericFlags := map[string]bool{
		"-k": true, "--kill-after": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: signalFlags, Placeholder: ""}, // preserve signal name verbatim
		{Flags: numericFlags, Placeholder: "N"},
	}

	var result []string
	parsingFlags := true
	durationSeen := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if parsingFlags && !durationSeen {
				durationSeen = true
			} else {
				parsingFlags = false
			}
			i++
			continue
		}

		if parsingFlags {
			// Check flag categories
			if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
				if placeholder == "" {
					// Signal flag: preserve the signal name verbatim
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
				} else {
					result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
				}
				continue
			}

			// Fused long flags (e.g. --kill-after=30s)
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}

			// Boolean flags
			if shellshape.IsFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}

			// First positional: DURATION
			if !durationSeen {
				if durationRE.MatchString(tok) {
					result = append(result, "N")
				} else {
					result = append(result, "N")
				}
				durationSeen = true
				i++
				continue
			}

			// Duration already seen, remaining is inner command
			parsingFlags = false
		}

		// Inner command tokens: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
