package shellshape

import "strings"

func init() {
	Register("lsblk", handleLsblk)
}

// handleLsblk handles the lsblk (list block devices) command.
// -o/--output consume a column list (<columns>).
// -e/--exclude and -I/--include consume major device numbers (N for single, <columns> for csv).
// -x/--sort consume a column name (<column>).
// -w/--width consume a numeric width (N).
// All other flags are boolean. Positionals are device paths via classifyToken.
func handleLsblk(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	columnListFlags := map[string]bool{
		"-o": true, "--output": true,
	}
	majorListFlags := map[string]bool{
		"-e": true, "--exclude": true,
		"-I": true, "--include": true,
	}
	sortFlags := map[string]bool{
		"-x": true, "--sort": true,
	}
	numericFlags := map[string]bool{
		"-w": true, "--width": true,
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

		if columnListFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<columns>")
				i++
			}
			continue
		}

		if majorListFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				val := args[i]
				if strings.Contains(val, ",") {
					result = append(result, "<columns>")
				} else {
					result = append(result, "N")
				}
				i++
			}
			continue
		}

		if sortFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<column>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional: device path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
