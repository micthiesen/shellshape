package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"npx", "bunx", "pnpx"} {
		shellshape.Register(name, handleNpx)
	}
}

// handleNpx handles npx and bunx (npm/bun package runners).
// Flags before the command: --yes/-y, --no (boolean), --package/-p (pkg name),
// --call/-c (code string), --workspace/-w (workspace name).
// The first non-flag positional is the package/command name (kept verbatim).
// All subsequent positionals are collapsed via classifyToken.
func handleNpx(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pkgFlags := map[string]bool{"-p": true, "--package": true}
	codeFlags := map[string]bool{"-c": true, "--call": true}
	valFlags := map[string]bool{"-w": true, "--workspace": true}

	var result []string
	commandSeen := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !commandSeen {
				commandSeen = true
			}
			i++
			continue
		}

		// Before the command name, handle flags that consume next token.
		if !commandSeen {
			if pkgFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<pkg>")
				continue
			}

			if codeFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<code>")
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

			// First positional: the command/package name, kept verbatim.
			result = append(result, tok)
			commandSeen = true
			i++
			continue
		}

		// After the command name: flags stay, positionals get classified.
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		classified := shellshape.ClassifyToken(tok)
		if classified == tok {
			// classifyToken didn't recognize it; collapse to <val>.
			classified = "<val>"
		}
		result = append(result, classified)
		i++
	}

	result = append(result, redirects...)
	return result
}
