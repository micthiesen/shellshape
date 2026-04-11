package shellshape

import "strings"

func init() {
	Register("scp", handleScp)
}

// isScpRemote detects scp remote specs: user@host:path, host:path, or bare user@host:
func isScpRemote(tok string) bool {
	// Must contain a colon (but not ://, which is a URI)
	idx := strings.Index(tok, ":")
	if idx < 0 {
		return false
	}
	if idx+2 < len(tok) && tok[idx+1] == '/' && tok[idx+2] == '/' {
		return false
	}
	// The part before colon must look like a host or user@host
	host := tok[:idx]
	if host == "" {
		return false
	}
	return true
}

// handleScp handles the scp command.
// Flags that consume the next token get appropriate placeholders.
// Remote specs (user@host:path) collapse to <remote>.
// Local paths use classifyToken.
func handleScp(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-D": true, "-F": true, "-i": true, "-S": true,
	}
	numericFlags := map[string]bool{
		"-l": true, "-P": true,
	}
	valFlags := map[string]bool{
		"-c": true, "-o": true, "-X": true,
	}
	hostFlags := map[string]bool{
		"-J": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{numericFlags, "N"},
		{valFlags, "<val>"},
		{hostFlags, "<host>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: remote spec or local path
		if isScpRemote(tok) {
			result = append(result, "<remote>")
		} else {
			result = append(result, classifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
