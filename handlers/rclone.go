package handlers

import (
	"strings"

	shellshape "github.com/micthiesen/shellshape"
)

func init() {
	shellshape.Register("rclone", handleRclone, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleRclone handles the rclone command.
// Remote paths (remote:path, :backend:path) and local paths → <path>.
// Filter flags (--filter, --include, --exclude) → <pattern>.
// Numeric flags (--transfers, --checkers) → N.
// Path flags (--log-file, --config) → <path>.
// Value flags (--bwlimit) → <val>.
func handleRclone(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	patternFlags := map[string]bool{
		"--filter":       true,
		"--include":      true,
		"--exclude":      true,
		"--include-from": true,
		"--exclude-from": true,
		"--filter-from":  true,
		"--files-from":   true,
	}

	numericFlags := map[string]bool{
		"--transfers":            true,
		"--checkers":             true,
		"--max-depth":            true,
		"--min-age":              true,
		"--max-age":              true,
		"--retries":              true,
		"--low-level-retries":    true,
		"--multi-thread-streams": true,
		"--buffer-size":          true,
	}

	pathFlags := map[string]bool{
		"--log-file":      true,
		"--config":        true,
		"--cache-dir":     true,
		"--temp-dir":      true,
		"--backup-dir":    true,
		"--password-file": true,
	}

	valFlags := map[string]bool{
		"--bwlimit":      true,
		"--log-level":    true,
		"--max-size":     true,
		"--min-size":     true,
		"--max-transfer": true,
		"--user-agent":   true,
		"--header":       true,
		"--rc-addr":      true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: patternFlags, Placeholder: "<pattern>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused --flag=value
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: remote paths and local paths → <path>
		result = append(result, rcloneClassifyPositional(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func rcloneClassifyPositional(tok string) string {
	// Remote path: remote:, remote:path, :backend:path
	if isRcloneRemotePath(tok) {
		return "<path>"
	}
	classified := shellshape.ClassifyToken(tok)
	if classified == "<path>" || strings.HasSuffix(classified, "-uri>") {
		return classified
	}
	return "<path>"
}

func isRcloneRemotePath(tok string) bool {
	// remote: or remote:path
	if idx := strings.IndexByte(tok, ':'); idx > 0 {
		prefix := tok[:idx]
		// Check it looks like a remote name (alphanumeric, no slashes before colon)
		for _, c := range prefix {
			if c == '/' {
				return false
			}
		}
		_ = prefix
		return true
	}
	// :backend:path
	if strings.HasPrefix(tok, ":") && len(tok) > 1 {
		return true
	}
	return false
}
