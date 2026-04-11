package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("file", handleFile)
}

// handleFile handles the file command (determine file type).
// Flags -m/-M/--magic-file, -f/--files-from consume a path argument.
// Flags -F/--separator, -e/--exclude, -P consume a value argument.
// All positional arguments are file paths.
func handleFile(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-m": true, "--magic-file": true,
		"-M": true,
		"-f": true, "--files-from": true,
	}
	valueFlags := map[string]bool{
		"-F": true, "--separator": true,
		"-e": true, "--exclude": true,
		"-P": true,
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

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
