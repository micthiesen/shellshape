package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"echo", "printf"} {
		shellshape.Register(name, handleEcho)
	}
}

// Real echo flags. Everything else that looks flag-ish (including
// dash-wrapped strings like `---chunkReload---`) is content.
var echoFlags = map[string]bool{
	"-n": true, "-e": true, "-E": true,
}

// handleEcho handles echo and printf commands.
// Only the canonical echo flags (`-n`, `-e`, `-E`) are preserved; any other
// dash-prefixed token is treated as content. All positional string args
// collapse to a single <str>. Subshells are kept verbatim.
func handleEcho(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	hasPositional := false

	for _, tok := range args {
		if echoFlags[tok] && !hasPositional {
			result = append(result, tok)
			continue
		}
		if shellshape.IsSubshellToken(tok) {
			if hasPositional {
				result = append(result, "<str>")
				hasPositional = false
			}
			result = append(result, tok)
			continue
		}
		hasPositional = true
	}
	if hasPositional {
		result = append(result, "<str>")
	}

	result = append(result, redirects...)
	return result
}
