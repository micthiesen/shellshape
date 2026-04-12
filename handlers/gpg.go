package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	for _, name := range []string{"gpg", "gpg2"} {
		shellshape.Register(name, handleGpg)
	}
}

// handleGpg handles gpg and gpg2.
// Recipient, keyserver, default-key flags collapse their arg to <val>.
// Output, homedir flags collapse their arg to <path>.
// All positionals collapse: paths and dotted filenames stay <path>,
// everything else becomes <val>.
func handleGpg(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
		"--homedir": true,
		"--keyring": true,
	}

	valFlags := map[string]bool{
		"-r": true, "--recipient": true, "--recipient-file": true,
		"--keyserver":   true,
		"--default-key": true,
		"--search-keys": true,
		"--local-user":  true, "-u": true,
		"--compress-algo": true, "--cipher-algo": true,
		"--digest-algo": true, "--cert-digest-algo": true,
		"--personal-cipher-preferences":   true,
		"--personal-digest-preferences":   true,
		"--personal-compress-preferences": true,
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

		// Positional: paths stay <path>, everything else becomes <val>
		result = append(result, gpgClassifyPositional(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func gpgClassifyPositional(tok string) string {
	classified := shellshape.ClassifyToken(tok)
	// Keep <path> and URL placeholders as-is
	if classified == "<path>" || strings.HasSuffix(classified, "-uri>") {
		return classified
	}
	// Dotted identifiers in gpg context are typically filenames (e.g. secret.txt.gpg, signature.asc)
	if classified == "<dotted-id>" {
		return "<path>"
	}
	// Everything else (key IDs, email addresses, bare words) → <val>
	return "<val>"
}
