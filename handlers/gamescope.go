package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("gamescope", handleGamescope)
}

// handleGamescope handles the gamescope Wayland compositor for gaming.
// Numeric flags (-w, -h, -W, -H, -r) consume the next arg as N.
// --backend consumes the next arg verbatim (structural).
// Boolean flags (--hdr-enabled, -f, -b, -e, --force-grab-cursor) pass through.
// After "--", all remaining tokens are the inner command, collapsed to generic
// placeholders (subshells are preserved).
func handleGamescope(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	numericFlags := map[string]bool{
		"-w": true, "-h": true, "-W": true, "-H": true, "-r": true,
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		// "--" separator: everything after is the inner command
		if tok == "--" {
			result = append(result, "--")
			i++
			for i < len(args) {
				inner := args[i]
				if shellshape.IsSubshellToken(inner) {
					result = append(result, inner)
				} else {
					result = append(result, collapseToVal(inner))
				}
				i++
			}
			break
		}

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if numericFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		// --backend takes a structural argument (keep verbatim)
		if tok == "--backend" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional before "--": collapse to generic placeholder
		result = append(result, collapseToVal(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// collapseToVal collapses a token to a generic placeholder.
// Paths are recognized as <path>, everything else becomes <val>.
func collapseToVal(tok string) string {
	classified := shellshape.ClassifyToken(tok)
	if classified != tok {
		return classified
	}
	return "<val>"
}
