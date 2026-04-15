package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("yabai", handleYabai)
}

// handleYabai handles the yabai macOS tiling window manager CLI.
//
// Top-level boolean management flags (--start-service, --restart-service,
// --load-sa, etc.) are kept verbatim.
//
// The -m flag introduces a domain (display, space, window, query, signal,
// config, rule, mouse, label) which is structural and kept verbatim. The
// next token after the domain is also structural: it is either a long-flag
// subcommand like --focus / --swap / --windows, or for "-m config" a config
// key name like window_border. Subsequent positional tokens are data
// (selectors like next/prev/recent, numeric IDs, config values) and collapse
// to <val> (or N when numeric). Subshells are always preserved verbatim.
func handleYabai(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	i := 0

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// -m <domain> <subcommand-or-key> [data...]
		if tok == "-m" {
			result = append(result, tok)
			i++
			// Domain token: kept verbatim.
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
					i++
					continue
				}
				result = append(result, args[i])
				i++
			}
			// Subcommand or config key: kept verbatim.
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
					i++
					continue
				}
				result = append(result, args[i])
				i++
			}
			// Remaining positionals for this -m group are data: collapse.
			// We stop at the next flag or the end so an unrelated trailing
			// flag (rare in practice) is still preserved structurally.
			for i < len(args) {
				next := args[i]
				if shellshape.IsSubshellToken(next) {
					result = append(result, next)
					i++
					continue
				}
				if shellshape.IsFlagToken(next) {
					break
				}
				if shellshape.NumberRE.MatchString(next) {
					result = append(result, "N")
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// Any other flag (top-level boolean management flags etc.): keep verbatim.
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Bare positional outside of an -m group: collapse generically.
		if shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
		} else {
			result = append(result, "<val>")
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
