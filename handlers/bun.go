package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	// Only these tokens are treated as bun subcommands. Anything else (a
	// script path like /tmp/foo.ts, or a package directory) is handed to
	// the handler as script-runner input.
	bunSubcommands := map[string]bool{
		"run": true, "test": true, "install": true, "i": true,
		"add": true, "remove": true, "rm": true, "uninstall": true,
		"update": true, "upgrade": true, "outdated": true,
		"build": true, "init": true, "create": true,
		"link": true, "unlink": true, "pm": true, "audit": true,
		"x": true, "repl": true, "patch": true,
	}
	shellshape.Register("bun", handleBun, shellshape.HandlerOptions{
		HasSubcommands: true,
		Subcommands:    bunSubcommands,
	})
}

// handleBun handles the bun command.
// The subcommand (run, test, install, build, x) is already extracted by
// normalizeSingleCommand. This handler processes the remaining tokens after
// the subcommand.
//
// Path flags (--env-file, --config, --outdir) consume the next token as <path>.
// Value flags (--target, --format) consume the next token as <val>.
// Numeric flags (--timeout, --bail, --port) consume the next token as N.
// Boolean flags (--watch, --hot, --frozen-lockfile, --production) are kept verbatim.
// Remaining positionals use classifyToken.
func handleBun(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"--env-file": true, "--config": true, "--outdir": true,
		"--outfile": true,
	}

	valueFlags := map[string]bool{
		"--target": true, "--format": true,
	}

	numericFlags := map[string]bool{
		"--timeout": true, "--bail": true, "--port": true,
	}

	// Script-runner mode: `bun <script.ts> [args...]`. Empty subcommand means
	// the caller invoked bun directly on a script path. First positional is a
	// path; remaining positionals are opaque script arguments.
	scriptMode := subcommand == ""

	var result []string
	sawScriptPath := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional
		if scriptMode {
			if !sawScriptPath {
				result = append(result, "<path>")
				sawScriptPath = true
			} else {
				result = append(result, "<arg>")
			}
			i++
			continue
		}
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
