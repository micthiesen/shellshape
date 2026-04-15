package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("find", handleFind)
}

// handleFind handles the find command.
// Pattern flags (-name, -iname, -path, etc.) make the next arg <pattern>
// unless it's a subshell (preserved verbatim).
// -exec/-execdir/-ok/-okdir consume everything until `;` or `+` and collapse
// the inner command to a single <cmd> placeholder.
// Positionals before the first primary (flag/expression) are treated as
// search roots and collapsed to <path>; bare directory names like `packages`
// would otherwise leak through.
// Subshells are preserved verbatim.
func handleFind(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	patternFlags := map[string]bool{
		"-name": true, "-iname": true,
		"-path": true, "-ipath": true,
		"-wholename": true, "-iwholename": true,
	}

	execFlags := map[string]bool{
		"-exec": true, "-execdir": true,
		"-ok": true, "-okdir": true,
	}

	var result []string
	// `seenPrimary` flips to true once we hit any find expression primary
	// (anything starting with `-` like `-type`, `-name`, or find's own
	// operators). Tokens before that are search roots.
	seenPrimary := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if execFlags[tok] {
			seenPrimary = true
			result = append(result, tok)
			i++
			// Consume the exec body until we hit `;` or `+`. Emit a single
			// <cmd> placeholder followed by the terminator.
			terminator := ""
			for i < len(args) {
				if args[i] == ";" || args[i] == "+" {
					terminator = args[i]
					i++
					break
				}
				i++
			}
			result = append(result, "<cmd>")
			if terminator != "" {
				result = append(result, terminator)
			}
			continue
		}

		if patternFlags[tok] {
			seenPrimary = true
			result = append(result, tok)
			i++
			if i < len(args) {
				next := args[i]
				if shellshape.IsSubshellToken(next) {
					result = append(result, next)
				} else {
					result = append(result, "<pattern>")
				}
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			seenPrimary = true
			result = append(result, tok)
			i++
			continue
		}

		// Tokens starting with - but not flags (find primaries like -delete, -print, -type, -o, -not).
		if strings.HasPrefix(tok, "-") {
			seenPrimary = true
			result = append(result, tok)
			i++
			continue
		}

		// Positional. Before any primary, these are search roots — collapse
		// to <path> even when they are bare directory names like `packages`.
		if !seenPrimary {
			result = append(result, "<path>")
			i++
			continue
		}

		// After a primary, positionals are values to primaries (e.g. the
		// argument to `-newer`, `-size`, or operator group terms). Fall back
		// to the generic classifier.
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
