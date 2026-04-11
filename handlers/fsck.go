package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("fsck", handleFsck)
}

// handleFsck handles the fsck (filesystem check) command.
// -t consumes the next token as <type> (filesystem type).
// All other flags are boolean. Positionals are device/filesystem paths
// collapsed via classifyToken.
func handleFsck(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -t/--type consumes the next token as a filesystem type.
		if tok == "-t" || tok == "--type" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<type>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device or filesystem path.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
