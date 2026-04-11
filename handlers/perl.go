package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("perl", handlePerl)
}

// handlePerl handles the perl command.
// Expression flags (-e, -E) make the next arg <perl-expr>.
// Bundled flags ending in 'e' or 'E' (like -ne, -pe, -lane) consume the next arg as <perl-expr>.
// Module flags (-M, -m) and fused forms (-MModule) collapse to -M <val> / -m <val>.
// Include flag (-I) and fused form (-I/path) collapse to -I <path>.
// First positional (if no expression assigned) becomes <script>.
// Remaining positionals use classifyToken.
func handlePerl(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	exprAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			exprAssigned = true
			i++
			continue
		}

		// Standalone -e / -E flag
		if tok == "-e" || tok == "-E" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<perl-expr>")
				}
				exprAssigned = true
				i++
			}
			continue
		}

		// Fused -M<module> or -m<module> (e.g. -MPOSIX, -mstrict)
		if (strings.HasPrefix(tok, "-M") || strings.HasPrefix(tok, "-m")) && len(tok) > 2 {
			result = append(result, tok[:2], "<val>")
			i++
			continue
		}

		// Fused -I<path> (e.g. -I/usr/lib)
		if strings.HasPrefix(tok, "-I") && len(tok) > 2 {
			result = append(result, "-I", "<path>")
			i++
			continue
		}

		// Separate -I flag
		if tok == "-I" {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		// Bundled short flags (e.g. -ne, -pe, -lane, -pi.bak, -0ne)
		// These start with - and contain only lowercase letters, digits, and possibly a dot-suffix.
		if isBundledPerlFlag(tok) {
			// Check if the bundle ends with 'e' or 'E' — meaning next arg is an expression.
			base := tok
			// Handle -pi.bak style: the .bak part is a fused -i extension
			if dotIdx := strings.IndexByte(tok[1:], '.'); dotIdx >= 0 {
				base = tok[:dotIdx+1]
			}
			lastChar := base[len(base)-1]
			if lastChar == 'e' || lastChar == 'E' {
				result = append(result, tok)
				i++
				if i < len(args) {
					if shellshape.IsSubshellToken(args[i]) {
						result = append(result, args[i])
					} else {
						result = append(result, "<perl-expr>")
					}
					exprAssigned = true
					i++
				}
				continue
			}
			// Bundled flags without trailing e/E — just pass through
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if !exprAssigned {
			result = append(result, "<script>")
			exprAssigned = true
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// isBundledPerlFlag returns true if tok looks like a bundled perl flag cluster
// (e.g. -ne, -pe, -lane, -pi.bak, -0ne). It must start with - and contain
// valid perl single-char flags. We use a permissive check: starts with -,
// length > 2, and doesn't start with -- (long flag).
func isBundledPerlFlag(tok string) bool {
	if len(tok) < 3 || !strings.HasPrefix(tok, "-") || strings.HasPrefix(tok, "--") {
		return false
	}
	// Check each character after the dash is a valid perl flag letter/digit or dot (for -i.bak)
	for j := 1; j < len(tok); j++ {
		c := tok[j]
		if c == '.' {
			// Everything after dot is the -i extension suffix, that's fine
			return true
		}
		if !isPerlFlagChar(c) {
			return false
		}
	}
	return true
}

func isPerlFlagChar(c byte) bool {
	// Valid perl single-character Flags: letters and digits used as flags
	// a c d e E h i l m n p s t T u U v w W x 0
	switch c {
	case 'a', 'c', 'd', 'e', 'E', 'h', 'i', 'l', 'n', 'p', 's', 't', 'T', 'u', 'U', 'v', 'w', 'W', 'x', '0':
		return true
	}
	return false
}
