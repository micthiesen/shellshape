package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"http", "https"} {
		shellshape.Register(name, handleHTTPie)
	}
}

// isHTTPMethod returns true if tok is an HTTP method keyword that HTTPie recognizes.
func isHTTPMethod(tok string) bool {
	switch tok {
	case "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS":
		return true
	}
	return false
}

// isHTTPieRequestItem returns true if tok looks like an HTTPie request item:
//
//	key=value   (data field)
//	key:value   (header)
//	key:=value  (raw JSON)
//	key==value  (query param)
//	key@file    (file upload)
//	key=@file   (file embed)
//
// We detect these by looking for the separator characters in a non-flag token.
func isHTTPieRequestItem(tok string) bool {
	// Must not be a flag or URL
	if strings.HasPrefix(tok, "-") {
		return false
	}
	if strings.Contains(tok, "://") {
		return false
	}

	// Check for HTTPie item separators: :=, ==, =@, @, :, =
	// Order matters: check multi-char separators first.
	for _, sep := range []string{":=@", ":=", "==", "=@", "@", "=", ":"} {
		idx := strings.Index(tok, sep)
		if idx > 0 {
			// Key part must not be empty and should look like a name (not a path or number).
			return true
		}
	}
	return false
}

// handleHTTPie handles the http/https commands (HTTPie).
// METHOD (GET/POST/PUT/DELETE/PATCH/HEAD/OPTIONS) is preserved verbatim.
// URL positionals are classified via classifyToken.
// Request items (key=val, key:val, key:=json, key==param, key@file) collapse to <item>.
// Auth flags (-a, --auth) collapse the next arg to <data>.
// Output/path flags (-o, --output, --cert, --cert-key, --session, --session-read-only) collapse to <path>.
// Numeric flags (--timeout, --max-redirects, --max-headers) collapse to N.
// Other arg-consuming flags collapse to <val>.
func handleHTTPie(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	dataFlags := map[string]bool{
		"-a": true, "--auth": true,
	}

	pathFlags := map[string]bool{
		"-o": true, "--output": true,
		"--cert": true, "--cert-key": true,
		"--session": true, "--session-read-only": true,
	}

	numericFlags := map[string]bool{
		"--timeout": true, "--max-redirects": true, "--max-headers": true,
	}

	valFlags := map[string]bool{
		"--auth-type": true,
		"--proxy":     true,
		"--verify":    true,
		"--print":     true, "-p": true,
		"--pretty": true,
		"--style":  true, "-s": true,
		"--format-options":   true,
		"--boundary":         true,
		"--response-charset": true,
		"--response-mime":    true,
		"--default-scheme":   true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: dataFlags, Placeholder: "<data>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		// Check for fused --flag=value syntax
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

		// HTTP method: preserve verbatim
		if isHTTPMethod(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Request item: collapse to <item>
		if isHTTPieRequestItem(tok) {
			result = append(result, "<item>")
			i++
			continue
		}

		// Positional: classify (URLs become <http-uri>, etc.)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
