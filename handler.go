package shellshape

import "strings"

// handlerFunc processes tokens after the executable name.
// subcommand is the subcommand token (empty if none).
type HandlerFunc func(subcommand string, tokens []string) []string

// flagCategory pairs a set of flag names with the placeholder to use for their values.
type FlagCategory struct {
	Flags       map[string]bool
	Placeholder string
}

// consumeFlagArg appends a flag and its next-token value (collapsed to placeholder)
// to result. Subshell tokens are preserved verbatim.
func ConsumeFlagArg(tok string, args []string, i int, result []string, placeholder string) ([]string, int) {
	result = append(result, tok)
	i++
	if i < len(args) {
		if IsSubshellToken(args[i]) {
			result = append(result, args[i])
		} else {
			result = append(result, placeholder)
		}
		i++
	}
	return result, i
}

// matchFlagCategory checks tok against a list of flag categories and returns
// the matching placeholder. Returns ("", false) if no category matches.
func MatchFlagCategory(tok string, categories []FlagCategory) (string, bool) {
	for _, cat := range categories {
		if cat.Flags[tok] {
			return cat.Placeholder, true
		}
	}
	return "", false
}

// consumeFusedFlag handles --flag=value syntax against a list of flag categories.
// Returns the normalized "flag=<placeholder>" string and true if matched,
// or ("", false) if the token doesn't contain '=' or doesn't match any category.
func ConsumeFusedFlag(tok string, categories []FlagCategory) (string, bool) {
	eqIdx := strings.IndexByte(tok, '=')
	if eqIdx < 0 {
		return "", false
	}
	key := tok[:eqIdx]
	for _, cat := range categories {
		if cat.Flags[key] {
			return key + "=" + cat.Placeholder, true
		}
	}
	return "", false
}

// EmitPositional appends a placeholder to result for a positional token,
// except when tok is a subshell substitution (`$(...)`) — in that case
// the token is appended verbatim so the subshell is preserved. Handlers
// must never collapse a subshell token to a data placeholder.
func EmitPositional(result []string, tok, placeholder string) []string {
	if IsSubshellToken(tok) {
		return append(result, tok)
	}
	return append(result, placeholder)
}

// ConsumeUntil scans args starting at index i and returns the first
// terminator token it finds along with the index of the slot after it.
// If no terminator is found, it returns ("", len(args)). Callers
// typically skip (or collapse to a single placeholder) everything the
// helper walked past.
//
// Use for grammars like `find ... -exec <cmd> ... ;|+`, where the body
// of a subcommand runs until a well-known terminator.
func ConsumeUntil(args []string, i int, terminators map[string]bool) (string, int) {
	for j := i; j < len(args); j++ {
		if terminators[args[j]] {
			return args[j], j + 1
		}
	}
	return "", len(args)
}

// RepairLeadingFlagSubcommand repairs the case where the outer normalizer
// extracted a global flag (either `--flag` or `--flag=value`) as what it
// believed was the subcommand. The flag token is already appended to
// result; this helper consumes the flag's value if it takes one, walks
// past any additional leading flags and their values, and emits the
// real subcommand token when it finds one.
//
// Parameters:
//   - subcommand: the token the outer normalizer extracted as the
//     subcommand. When this is not flag-like, the helper is a no-op.
//   - args: the handler's remaining tokens.
//   - i: the current index in args.
//   - result: the handler's accumulating output slice.
//   - categories: flag categories used to decide which leading flags
//     consume a next-token value and which placeholder that value gets.
//     Pass nil when every leading flag is boolean (e.g. `ufw --dry-run`).
//   - verbatimValueFlags: optional flags whose values are kept verbatim
//     instead of being replaced with a placeholder. Useful for small
//     finite value sets like systemctl's `-t/--type`. Pass nil if unused.
//
// Returns the extended result slice, the real subcommand token (or "" if
// none was found), and the new position in args.
//
// Fused-flag note: when the normalizer extracted a `--flag=value` token,
// it was already appended to result verbatim and cannot be rewritten —
// the helper simply advances and finds the real subcommand.
func RepairLeadingFlagSubcommand(
	subcommand string,
	args []string,
	i int,
	result []string,
	categories []FlagCategory,
	verbatimValueFlags map[string]bool,
) ([]string, string, int) {
	if subcommand == "" || !looksLikeFlagSubcommand(subcommand) {
		return result, subcommand, i
	}

	// 1. Consume the leaked flag's value, if it has one. A fused flag
	//    (`--flag=value`) already carries its value in the token that was
	//    appended by the normalizer; nothing to consume.
	if !strings.ContainsRune(subcommand, '=') {
		if placeholder, ok := MatchFlagCategory(subcommand, categories); ok {
			if i < len(args) && !IsFlagToken(args[i]) {
				result = EmitPositional(result, args[i], placeholder)
				i++
			}
		} else if verbatimValueFlags[subcommand] {
			if i < len(args) && !IsFlagToken(args[i]) {
				result = append(result, args[i])
				i++
			}
		}
	}

	// 2. Walk forward through any further leading flags (with their
	//    values) and subshells until we hit a real subcommand token.
	for i < len(args) {
		tok := args[i]

		if IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if fused, ok := ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if placeholder, ok := MatchFlagCategory(tok, categories); ok {
			result, i = ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if verbatimValueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) && !IsFlagToken(args[i]) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		if IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First non-flag positional: the real subcommand.
		result = append(result, tok)
		return result, tok, i + 1
	}

	return result, "", i
}

// looksLikeFlagSubcommand reports whether a subcommand token came from
// the outer normalizer extracting a flag-like thing. This covers plain
// flags (`-t`, `--filter`) and fused flags (`--flag=value`).
func looksLikeFlagSubcommand(tok string) bool {
	if IsFlagToken(tok) {
		return true
	}
	// A `--flag=value` token that begins with a dash is already handled
	// by IsFlagToken; this extra branch catches `-Xval` style fused
	// forms that IsFlagToken may not recognize.
	if strings.HasPrefix(tok, "-") && strings.ContainsRune(tok, '=') {
		return true
	}
	return false
}

var RedirectConsumeNext = map[string]bool{
	">": true, ">>": true, "<": true,
	"&>": true, "&>>": true,
	"2>": true, "2>>": true, "1>": true,
	"<<<": true,
}

var RedirectStandalone = map[string]bool{
	"2>&1": true, "1>&2": true,
	"&>/dev/null": true, "2>/dev/null": true,
}

// splitRedirects partitions tokens into non-redirect args and redirect tail.
func SplitRedirects(tokens []string) (args, redirects []string) {
	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		if RedirectStandalone[tok] {
			redirects = append(redirects, tok)
			i++
			continue
		}

		if RedirectConsumeNext[tok] {
			if i+1 < len(tokens) {
				placeholder := "<path>"
				if tok == "<<<" {
					placeholder = "<str>"
				}
				redirects = append(redirects, tok, placeholder)
				i += 2
			} else {
				redirects = append(redirects, tok)
				i++
			}
			continue
		}

		args = append(args, tok)
		i++
	}
	return
}
