package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	shellshape.Register("fzf", handleFzf)
}

func handleFzf(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-q": true, "--query": true,
		"-f": true, "--filter": true,
		"--prompt": true, "--header": true,
		"--preview": true, "--bind": true,
		"--height": true,
		"-d":       true, "--delimiter": true,
		"--with-nth": true, "--nth": true,
		"--preview-window": true,
		"--border-label":   true, "--header-label": true,
		"--color": true, "--info": true,
		"--pointer": true, "--marker": true,
		"--tabstop": true, "--hscroll-off": true,
		"--jump-labels": true,
	}
	structuralFlags := map[string]bool{
		"--layout": true,
	}
	categories := []shellshape.FlagCategory{
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

		// Fused --flag=value: structural flags keep value, others use category
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 {
			key := tok[:eqIdx]
			if structuralFlags[key] {
				result = append(result, tok)
				i++
				continue
			}
			if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, fused)
				i++
				continue
			}
		}

		// Structural flags: keep value verbatim
		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
