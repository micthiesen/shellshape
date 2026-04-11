package shellshape

import "strings"

// handleTar handles the tar command.
// Recognizes bundled flags (czf, xvzf) with or without a leading dash.
// If the bundle contains 'f', the next token is consumed as <path> (the archive file).
// Flags like -C, -T, -X, --file, --directory consume the next token as <path>.
// --exclude and --include consume the next token as <pattern>.
// --strip-components and --block-size consume the next token as N.
// Remaining positionals use classifyToken.
func handleTar(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next argument is the archive file (distinct placeholder).
	archiveFlags := map[string]bool{
		"-f": true, "--file": true,
	}

	// Flags whose next argument is a path.
	pathFlags := map[string]bool{
		"-C": true, "--cd": true, "--directory": true,
		"-T": true, "--files-from": true,
		"-X": true, "--exclude-from": true,
	}

	// Flags whose next argument is a pattern.
	patternFlags := map[string]bool{
		"--exclude": true, "--include": true,
	}

	// Flags whose next argument is numeric.
	numericFlags := map[string]bool{
		"--strip-components": true,
		"-b": true, "--block-size": true,
	}

	// Flags whose next argument is a generic value.
	valFlags := map[string]bool{
		"--format": true,
		"--owner": true, "--group": true,
		"--gname": true, "--uname": true,
		"--gid": true, "--uid": true,
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Bundled flags like czf, xvzf (with or without leading dash).
		// Detect: first token or a token that looks like a bundle of single letters.
		if i == 0 && !strings.HasPrefix(tok, "--") && isTarBundle(tok) {
			bundle := tok
			result = append(result, bundle)
			i++
			// If bundle contains 'f', next arg is the archive file.
			raw := strings.TrimPrefix(bundle, "-")
			if strings.Contains(raw, "f") && i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<archive>")
				}
				i++
			}
			continue
		}

		if archiveFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<archive>")
				}
				i++
			}
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

		if patternFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<pattern>")
				}
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file or directory path.
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// isTarBundle returns true if the token looks like a tar bundled-flag word.
// e.g., "czf", "xvzf", "-czf", "-xvzf". Must be all ASCII letters (after
// stripping optional leading dash) and contain at least one tar mode letter.
func isTarBundle(tok string) bool {
	raw := strings.TrimPrefix(tok, "-")
	if len(raw) < 2 {
		return false
	}
	for _, c := range raw {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}
	// Must contain at least one mode letter.
	for _, c := range raw {
		switch c {
		case 'c', 'x', 't', 'r', 'u':
			return true
		}
	}
	return false
}
