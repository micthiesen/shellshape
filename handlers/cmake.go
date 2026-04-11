package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("cmake", handleCmake)
}

// handleCmake handles cmake in all its modes:
// configure (default), --build, --install, -P (script), -E (tool).
//
// Fused -DVAR=VALUE flags are split into -D <val>.
// Generator (-G), toolset (-T), platform (-A) flags consume next arg as <val>.
// Source (-S), build (-B), cache (-C) flags consume next arg as <path>.
// In --build mode, --target/-t consume <val>, -j/--parallel consume N.
// In --install mode, --prefix consumes <path>, --component/--config consume <val>.
// Remaining positionals use classifyToken.
func handleCmake(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-G": true, "-T": true, "-A": true,
		"-U":                           true,
		"--preset":                     true,
		"--log-level":                  true,
		"--trace-format":               true,
		"--profiling-format":           true,
		"--resolve-package-references": true,
		"--config":                     true,
		"--target":                     true, "-t": true,
		"--component":                     true,
		"--default-directory-permissions": true,
		"--list-presets":                  true,
		"--debug-find-pkg":                true,
		"--debug-find-var":                true,
		"--format":                        true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-S": true, "-B": true, "-C": true,
		"-P":                 true,
		"--toolchain":        true,
		"--install-prefix":   true,
		"--prefix":           true,
		"--trace-source":     true,
		"--trace-redirect":   true,
		"--graphviz":         true,
		"--profiling-output": true,
		"--sarif-output":     true,
		"--debugger-pipe":    true,
		"--debugger-dap-log": true,
	}

	// Flags whose next token is a number.
	numericFlags := map[string]bool{
		"-j": true, "--parallel": true,
	}

	// Fused flags that use = syntax (for consumeFusedFlag).
	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		// Fused -DVAR=VALUE or -DVAR:TYPE=VALUE → -D <val>
		if strings.HasPrefix(tok, "-D") && len(tok) > 2 && tok[2] != '-' {
			result = append(result, "-D", "<val>")
			i++
			continue
		}

		// Separate -D flag: next token is the variable definition.
		if tok == "-D" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// Mode flags that consume the next token as a directory.
		// Use classifyToken so that "." stays verbatim.
		if tok == "--build" || tok == "--install" || tok == "--open" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, shellshape.ClassifyToken(args[i]))
				}
				i++
			}
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

		// Positional: classify generically
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
