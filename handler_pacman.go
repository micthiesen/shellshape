package shellshape

import "strings"

func init() {
	for _, name := range []string{"pacman", "yay", "paru"} {
		Register(name, handlePacman, HandlerOptions{HasSubcommands: false})
	}
}

// handlePacman handles pacman, yay, and paru (Arch Linux package managers).
//
// Pacman uses operation flags instead of subcommands: -S (sync), -R (remove),
// -Q (query), -U (upgrade from file), -D (database), -F (files). Modifiers
// are bundled with the operation flag (e.g. -Syu, -Ss, -Qi).
//
// Package names are structural and kept verbatim for install/remove/query.
// Search queries (-Ss, -Qs, -Fs) collapse to <query>.
// File paths (-U, -Qo) use classifyToken.
// Global flags like --config, --dbpath, --root consume the next token as a path or value.
func handlePacman(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags that consume the next token as a path.
	pathFlags := map[string]bool{
		"--config":   true,
		"--dbpath":   true,
		"-b":         true,
		"--root":     true,
		"-r":         true,
		"--cachedir": true,
		"--logfile":  true,
		"--gpgdir":   true,
		"--hookdir":  true,
	}

	// Flags that consume the next token as a generic value.
	valFlags := map[string]bool{
		"--arch":         true,
		"--color":        true,
		"--print-format": true,
	}

	// Long-form operation flags to their short letter.
	longOps := map[string]byte{
		"--sync":     'S',
		"--remove":   'R',
		"--query":    'Q',
		"--upgrade":  'U',
		"--database": 'D',
		"--files":    'F',
	}

	// Parse all flags to determine the operation mode and which modifiers are active.
	var operation byte // S, R, Q, U, D, F
	modifiers := map[byte]bool{}

	for _, tok := range args {
		if op, ok := longOps[tok]; ok {
			operation = op
			continue
		}
		// Bundled short flags like -Syu, -Qi, -Rcs
		if strings.HasPrefix(tok, "-") && !strings.HasPrefix(tok, "--") && len(tok) >= 2 {
			for j := 1; j < len(tok); j++ {
				c := tok[j]
				switch c {
				case 'S', 'R', 'Q', 'U', 'D', 'F':
					operation = c
				default:
					modifiers[c] = true
				}
			}
		}
	}

	// Determine if we're in a "search" mode where positionals are queries.
	isSearch := modifiers['s'] && (operation == 'S' || operation == 'Q' || operation == 'F')
	// -Qo: query which package owns a file (positional is a file path).
	isOwns := modifiers['o'] && operation == 'Q'
	// -U: upgrade from local file (positional is a file path).
	isFileUpgrade := operation == 'U'

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Positional argument handling depends on operation mode.
		if isSearch {
			result = append(result, "<query>")
			i++
			continue
		}

		if isOwns || isFileUpgrade {
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		// Default: package names are structural, keep verbatim.
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
