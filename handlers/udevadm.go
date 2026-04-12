package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("udevadm", handleUdevadm, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleUdevadm handles udevadm subcommand arguments.
// Subcommands: control, trigger, info, settle, monitor, test, hwdb.
// For info: -q/--query value is structural (property, name, etc.), -n/--name and -p/--path collapse to <path>.
// For trigger: --type and --action values are structural (kept verbatim).
// For control/monitor: mostly boolean flags.
func handleUdevadm(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose values are structural (kept verbatim)
	structuralFlags := map[string]bool{
		"-q": true, "--query": true,
		"--type":   true,
		"--action": true,
	}

	// Flags whose values collapse to <path>
	pathFlags := map[string]bool{
		"-n": true, "--name": true,
		"-p": true, "--path": true,
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

		// Handle --flag=value fused syntax
		if strings.Contains(tok, "=") && shellshape.IsFlagToken(tok) {
			eqIdx := strings.IndexByte(tok, '=')
			key := tok[:eqIdx]
			val := tok[eqIdx+1:]

			if pathFlags[key] {
				result = append(result, key+"=<path>")
			} else {
				// Structural: keep verbatim (--query=name, --type=subsystems, --action=add)
				_ = val
				result = append(result, tok)
			}
			i++
			continue
		}

		if structuralFlags[tok] {
			// Keep next token verbatim
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional - classify generically (typically device/sys paths)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
