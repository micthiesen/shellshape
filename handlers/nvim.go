package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
	"strings"
)

func init() {
	for _, name := range []string{"nvim", "vim", "vi"} {
		shellshape.Register(name, handleNvim)
	}
}

var plusLineRE = regexp.MustCompile(`^\+\d+$`)

// handleNvim handles nvim, vim, vi.
// -c <command> → <val>, -u <vimrc> → <path>, -S <file> → <path>.
// +N (line number) → +N. +<anything else> → +<val>.
// Boolean flags: -d, -R, -p, -o, -O, --headless, --clean, etc.
// Positional args are file paths → ClassifyToken.
func handleNvim(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{"-c": true}
	pathFlags := map[string]bool{"-u": true, "-S": true}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// +N line number
		if plusLineRE.MatchString(tok) {
			result = append(result, "+N")
			i++
			continue
		}

		// +command or +/pattern
		if strings.HasPrefix(tok, "+") && len(tok) > 1 {
			result = append(result, "+<val>")
			i++
			continue
		}

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
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

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
