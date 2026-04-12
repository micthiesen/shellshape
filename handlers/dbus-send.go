package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("dbus-send", handleDbusSend)
}

// D-Bus type prefixes that precede a value argument.
var dbusTypePrefixes = []string{
	"string:", "int16:", "uint16:", "int32:", "uint32:",
	"int64:", "uint64:", "byte:", "boolean:", "double:",
	"objpath:", "array:", "dict:", "variant:",
}

// Flags whose =value should become =<val>.
var dbusValFlags = map[string]bool{
	"--dest": true, "--type": true,
}

// Flags whose =value should become =N.
var dbusNumericFlags = map[string]bool{
	"--reply-timeout": true,
}

// Flags whose =value should be kept verbatim (small fixed set of literals).
var dbusLiteralFlags = map[string]bool{
	"--print-reply": true,
}

// handleDbusSend normalizes dbus-send commands.
// Flags like --dest=, --type= collapse their value; --reply-timeout= becomes N.
// --print-reply and --print-reply=literal are kept verbatim.
// Positional D-Bus typed arguments (string:foo, int32:42) become TYPE:<val>.
// Other positionals use ClassifyToken.
func handleDbusSend(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string

	for i := 0; i < len(args); i++ {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// Handle --flag=value fused flags.
		if strings.HasPrefix(tok, "--") && strings.Contains(tok, "=") {
			eqIdx := strings.IndexByte(tok, '=')
			key := tok[:eqIdx]
			if dbusValFlags[key] {
				result = append(result, key+"=<val>")
				continue
			}
			if dbusNumericFlags[key] {
				result = append(result, key+"=N")
				continue
			}
			if dbusLiteralFlags[key] {
				// Keep verbatim (e.g. --print-reply=literal).
				result = append(result, tok)
				continue
			}
			// Unknown fused flag: collapse value.
			result = append(result, key+"=<val>")
			continue
		}

		// Boolean flags (no =).
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// D-Bus typed argument: TYPE:VALUE -> TYPE:<val>.
		if prefix := dbusTypePrefix(tok); prefix != "" {
			result = append(result, prefix+"<val>")
			continue
		}

		// Regular positional.
		result = append(result, shellshape.ClassifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}

// dbusTypePrefix returns the "type:" prefix if tok starts with a known D-Bus
// type prefix, or "" otherwise.
func dbusTypePrefix(tok string) string {
	for _, p := range dbusTypePrefixes {
		if strings.HasPrefix(tok, p) && len(tok) > len(p) {
			return p
		}
	}
	return ""
}
