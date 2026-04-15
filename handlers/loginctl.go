package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

func init() {
	shellshape.Register("loginctl", handleLoginctl, shellshape.HandlerOptions{HasSubcommands: true})
}

var loginctlNumericRE = regexp.MustCompile(`^\d+$`)

// Session-taking subcommands where positional is a session ID (numeric → N).
var loginctlSessionSubcmds = map[string]bool{
	"show-session": true, "activate": true,
	"lock-session": true, "unlock-session": true,
	"terminate-session": true,
}

// handleLoginctl handles loginctl subcommand arguments.
// Session IDs → N, usernames and other values → <val>.
func handleLoginctl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-H": true, "--host": true,
		"-M": true, "--machine": true,
		"-p": true, "--property": true,
		"--kill-who": true,
		"--signal":   true,
		"-s":         true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0

	result, subcommand, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, categories, nil,
	)

	isSessionSubcmd := loginctlSessionSubcmds[subcommand]
	hasPositional := false

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

		// Positional
		if !hasPositional {
			if isSessionSubcmd && loginctlNumericRE.MatchString(tok) {
				result = append(result, "N")
			} else {
				result = append(result, "<val>")
			}
			hasPositional = true
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
