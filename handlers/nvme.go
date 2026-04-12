package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("nvme", handleNvme, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleNvme handles nvme-cli NVMe device management commands.
// Subcommands (list, smart-log, id-ctrl, etc.) are consumed by the framework.
// Numeric flags: -n/--namespace-id, -l/--lba-format, -s/--ses, -f/--feature-id,
// -V/--value, -a/--action, -i/--log-id, etc. → N
// Structural flags: -o/--output-format (json/normal) — kept verbatim.
// Boolean flags: -H/--human-readable, -v/--verbose, etc.
// Positional is always a device path → <path>.
func handleNvme(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-n": true, "--namespace-id": true,
		"-l": true, "--lba-format": true,
		"-s": true, "--ses": true,
		"-f": true, "--feature-id": true,
		"-V": true, "--value": true,
		"-a": true, "--action": true,
		"-i": true, "--log-id": true,
		"-c": true, "--cdw10": true,
		"--cdw11": true, "--cdw12": true,
		"--cdw13": true, "--cdw14": true, "--cdw15": true,
		"-d": true, "--data-len": true,
	}

	// Structural flags: keep value verbatim (defines output mode, not data)
	structuralFlags := map[string]bool{
		"-o": true, "--output-format": true,
	}

	// Long flags that keep value verbatim in fused form
	structuralLongFlags := map[string]bool{
		"--output-format": true,
	}

	numericLongFlags := map[string]bool{
		"--namespace-id": true, "--lba-format": true, "--ses": true,
		"--feature-id": true, "--value": true, "--action": true,
		"--log-id": true, "--data-len": true,
		"--cdw10": true, "--cdw11": true, "--cdw12": true,
		"--cdw13": true, "--cdw14": true, "--cdw15": true,
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

		// Handle --flag=value
		if eqIdx := strings.IndexByte(tok, '='); eqIdx > 0 && strings.HasPrefix(tok, "--") {
			key := tok[:eqIdx]
			if structuralLongFlags[key] {
				result = append(result, tok)
				i++
				continue
			}
			if numericLongFlags[key] {
				result = append(result, key+"=N")
				i++
				continue
			}
		}

		// Structural flags: keep next token verbatim
		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Numeric flags
		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: device path or namespace number
		if shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
