package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	shellshape.Register("lms", handleLms, shellshape.HandlerOptions{HasSubcommands: true})
}

// Known lms sub-subcommands that should stay verbatim.
var lmsSubSubcommands = map[string]bool{
	"start": true, "stop": true, "status": true,
}

// handleLms handles the LM Studio CLI.
// Sub-subcommands (start, stop, status) stay verbatim.
// Model names and other positionals → <val>.
// --port, --context-length → N.
// --gpu → <val>.
// --gguf and path flags → <path>.
func handleLms(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"--gguf":       true,
		"--model-path": true,
	}

	numericFlags := map[string]bool{
		"--port":           true,
		"--context-length": true,
		"--threads":        true,
		"--batch-size":     true,
	}

	valFlags := map[string]bool{
		"--gpu":    true,
		"--preset": true,
		"--format": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	subSubSeen := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused --flag=value
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

		// Sub-subcommand: keep verbatim
		if !subSubSeen && lmsSubSubcommands[tok] {
			result = append(result, tok)
			subSubSeen = true
			i++
			continue
		}

		// Positional: paths stay <path>, model names and everything else → <val>
		result = append(result, lmsClassifyPositional(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func lmsClassifyPositional(tok string) string {
	classified := shellshape.ClassifyToken(tok)
	if classified == "<path>" || strings.HasSuffix(classified, "-uri>") {
		return classified
	}
	return "<val>"
}
