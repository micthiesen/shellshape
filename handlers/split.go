package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("split", handleSplit)
	shellshape.Register("csplit", handleCsplit)
}

// handleSplit handles the split command.
// -l (line count) and -n (chunk count) consume the next token as N.
// -b (byte count, e.g. 10M) consumes the next token as <size>.
// -a (suffix length) consumes the next token as N.
// -p (pattern) consumes the next token as <pattern>.
// -c, -d are boolean flags.
// First positional is the input file (classifyToken), second is the output prefix (<prefix>).
func handleSplit(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{"-l": true, "-n": true, "-a": true}

	var result []string
	positionalIndex := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			positionalIndex++
			continue
		}

		// Flags that consume a numeric argument
		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// -b consumes a size argument (e.g. 10M, 512k)
		if tok == "-b" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<size>")
			continue
		}

		// -p consumes a pattern argument
		if tok == "-p" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<pattern>")
			continue
		}

		// Other flags (boolean or long-form)
		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positionals: first is file, second is prefix
		if positionalIndex == 0 {
			result = append(result, shellshape.ClassifyToken(tok))
		} else {
			result = append(result, "<prefix>")
		}
		positionalIndex++
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleCsplit handles the csplit command.
// -f (prefix) consumes the next token as <prefix>.
// -n (digit count) consumes the next token as N.
// -b (suffix format, GNU) consumes the next token as <format>.
// -k, -s are boolean flags.
// First positional is the input file (classifyToken).
// Remaining positionals are split arguments (patterns, line numbers, repeat specs) → <split-arg>.
func handleCsplit(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	positionalIndex := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			positionalIndex++
			continue
		}

		// -f consumes a prefix argument
		if tok == "-f" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<prefix>")
			continue
		}

		// -n consumes a numeric argument
		if tok == "-n" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// -b consumes a suffix format argument (GNU extension)
		if tok == "-b" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<format>")
			continue
		}

		// Other flags (boolean: -k, -s, or long-form)
		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positionals
		if positionalIndex == 0 {
			// First positional: input file (or "-" for stdin)
			if tok == "-" {
				result = append(result, "-")
			} else {
				result = append(result, shellshape.ClassifyToken(tok))
			}
		} else {
			// Remaining positionals: split arguments (patterns, line numbers, repeat specs)
			result = append(result, "<split-arg>")
		}
		positionalIndex++
		i++
	}

	result = append(result, redirects...)
	return result
}
