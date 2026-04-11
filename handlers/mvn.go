package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("mvn", handleMvn)
}

// handleMvn handles Apache Maven commands.
//
// Phases and plugin goals (clean, compile, package, exec:java, dependency:tree)
// are structural and kept verbatim.
//
// -D properties: -DskipTests (no =) stays verbatim; -Dkey=value collapses value
// to key=<val>. Separate form (-D key) keeps the key verbatim.
//
// -P profiles: fused -Pname splits to -P <val>; separate -P name collapses to
// -P <val>.
//
// Path flags (-f, -s, -gs, -l) collapse next token to <path>.
// Value flags (-pl, -rf, --define, etc.) collapse next token to <val>.
// Numeric flags (-T) collapse next token to N.
// Boolean flags (-B, -U, -X, -q, -e, -o, etc.) are kept verbatim.
func handleMvn(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"-s": true, "--settings": true,
		"-gs": true, "--global-settings": true,
		"-l": true, "--log-file": true,
	}

	valFlags := map[string]bool{
		"-pl": true, "--projects": true,
		"-rf": true, "--resume-from": true,
		"-am": true, "--also-make": true,
		"-amd": true, "--also-make-dependents": true,
		"--define": true,
	}

	numFlags := map[string]bool{
		"-T": true, "--threads": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: numFlags, Placeholder: "N"},
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

		// Handle -Dkey=value and -Dkey (fused system properties)
		if strings.HasPrefix(tok, "-D") && len(tok) > 2 {
			rest := tok[2:]
			if eqIdx := strings.IndexByte(rest, '='); eqIdx >= 0 {
				// -Dkey=value -> -Dkey=<val>
				key := rest[:eqIdx]
				result = append(result, "-D"+key+"=<val>")
			} else {
				// -DskipTests (boolean property, no value)
				result = append(result, tok)
			}
			i++
			continue
		}

		// Handle -Pprofile (fused profile)
		if strings.HasPrefix(tok, "-P") && len(tok) > 2 {
			result = append(result, "-P", "<val>")
			i++
			continue
		}

		// Handle -P as separate flag (consume next token)
		if tok == "-P" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Handle -D as separate flag (next token is the property key, keep it verbatim)
		if tok == "-D" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Standard flag categories
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Fused long flags with = (e.g. --define=value)
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

		// Positionals: phases and plugin goals are structural, keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
