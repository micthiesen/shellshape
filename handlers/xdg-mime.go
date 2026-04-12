package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("xdg-mime", handleXdgMime, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleXdgMime handles xdg-mime subcommand arguments.
// Subcommands:
//   - query default <mime-type>: MIME type → <val>
//   - query filetype <path>: path via ClassifyToken
//   - default <app.desktop> <mime-type>...: all positionals → <val>
//   - install/uninstall [--mode <val>] [--novendor] <path>: path via ClassifyToken
//
// The --mode flag consumes a value argument.
func handleXdgMime(subcommand string, tokens []string) []string {
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

		// Handle --mode flag (takes an argument: "user" or "system")
		if tok == "--mode" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional handling depends on subcommand
		switch subcommand {
		case "query":
			// First positional is the sub-subcommand (default/filetype), keep verbatim.
			// Second positional: for "filetype" it's a path, for "default" it's a MIME type.
			if len(result) == 0 {
				// Sub-subcommand: keep verbatim
				result = append(result, tok)
			} else if len(result) > 0 && result[0] == "filetype" {
				result = append(result, shellshape.ClassifyToken(tok))
			} else {
				// "default" sub-subcommand: MIME type → <val>
				result = append(result, "<val>")
			}

		case "default":
			// All positionals are .desktop name or MIME types → <val>
			result = append(result, "<val>")

		case "install", "uninstall":
			// Positionals are file paths
			result = append(result, shellshape.ClassifyToken(tok))

		default:
			result = append(result, shellshape.ClassifyToken(tok))
		}

		i++
	}

	result = append(result, redirects...)
	return result
}
