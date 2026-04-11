package shellshape

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	headVariantRE = regexp.MustCompile(`^HEAD(?:[~^][~^\d]*|@\{[^}]*\})$`)

	pathExtensionsRE = regexp.MustCompile(`(?i)\.(?:py|js|ts|tsx|jsx|go|rs|rb|java|kt|swift|c|cc|cpp|h|hpp|md|mdx|txt|json|jsonc|yaml|yml|toml|ini|cfg|conf|sh|bash|zsh|fish|sql|html|htm|css|scss|sass|less|lock|log|csv|tsv|xml|plist|png|jpg|jpeg|gif|svg|webp|ico|pdf|zip|tar|gz|bz2|xz|7z|rar|dmg|env|gitignore|dockerignore|editorconfig|prettierrc|eslintrc|pem|crt|csr|der|p12|pfx|key|cer)$`)

	uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

	hexRE = regexp.MustCompile(`^[0-9a-fA-F]+$`)

	numberRE = regexp.MustCompile(`^-?\d+(\.\d+)?$`)

	schemeRE = regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9+.-]*)://`)

	gitRemoteRE = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+:.+$`)

	revRangeRE = regexp.MustCompile(`^(.+?)(\.\.\.?)(.+)$`)

	dottedIDRE = regexp.MustCompile(`^[a-zA-Z_]\w*(\.[a-zA-Z_]\w*)+$`)

	envAssignRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)
)

func classifyToken(tok string) string {
	// 1. Long flag
	if strings.HasPrefix(tok, "--") && len(tok) > 2 {
		if idx := strings.Index(tok, "="); idx >= 0 {
			return tok[:idx] + "=<val>"
		}
		return tok
	}

	// 2. Short flag: starts with - and first char after dashes is alpha
	if strings.HasPrefix(tok, "-") && len(tok) >= 2 {
		rest := strings.TrimLeft(tok[1:], "-")
		if len(rest) > 0 && unicode.IsLetter(rune(rest[0])) {
			return tok
		}
		// -10 (all digits after dash) falls through to number check
	}

	// 3. Signed number
	if numberRE.MatchString(tok) {
		return "N"
	}

	// 4. URL with scheme
	if m := schemeRE.FindStringSubmatch(tok); m != nil {
		scheme := strings.ToLower(m[1])
		return "<" + scheme + "-uri>"
	}

	// 5. SSH/git remote
	if gitRemoteRE.MatchString(tok) {
		return "<git-uri>"
	}

	// 6. UUID
	if uuidRE.MatchString(tok) {
		return "<uuid>"
	}

	// 7. Git SHA hex 7-40 chars, not all digits
	if len(tok) >= 7 && len(tok) <= 40 && hexRE.MatchString(tok) {
		allDigits := true
		for _, c := range tok {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if !allDigits {
			return "<hash>"
		}
	}

	// 8. Rev range (before path, since ranges can contain slashes)
	if m := revRangeRE.FindStringSubmatch(tok); m != nil {
		if looksLikeRev(m[1]) && looksLikeRev(m[3]) {
			return "<range>"
		}
	}

	// 9. HEAD variants (HEAD~5, HEAD^, HEAD^^, HEAD@{upstream}) but bare HEAD stays
	if headVariantRE.MatchString(tok) {
		return "<rev>"
	}

	// 10. Path: starts with /, ~, ./, ../, or contains /
	if strings.HasPrefix(tok, "/") || strings.HasPrefix(tok, "~") ||
		strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") ||
		strings.Contains(tok, "/") {
		return "<path>"
	}

	// 11. Known file extension
	if pathExtensionsRE.MatchString(tok) {
		return "<path>"
	}

	// 12. Dotted identifier (2+ segments, each starts with letter/_)
	if dottedIDRE.MatchString(tok) {
		return "<dotted-id>"
	}

	// 13. Keep verbatim
	return tok
}

func looksLikeRev(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s == "HEAD" {
		return true
	}
	if headVariantRE.MatchString(s) {
		return true
	}
	// 7-40 hex chars
	if len(s) >= 7 && len(s) <= 40 && hexRE.MatchString(s) {
		return true
	}
	// Branch/ref name: starts with letter/digit, may contain /.-_
	if unicode.IsLetter(rune(s[0])) || unicode.IsDigit(rune(s[0])) {
		for _, c := range s {
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '/' && c != '.' && c != '-' && c != '_' {
				return false
			}
		}
		return true
	}
	return false
}

func isSubshellToken(tok string) bool {
	return strings.HasPrefix(tok, "$(") && strings.HasSuffix(tok, ")")
}

func isFlagToken(tok string) bool {
	if !strings.HasPrefix(tok, "-") || len(tok) < 2 {
		return false
	}
	rest := strings.TrimLeft(tok[1:], "-")
	return len(rest) > 0 && unicode.IsLetter(rune(rest[0]))
}

func isEnvAssignment(tok string) bool {
	return envAssignRE.MatchString(tok)
}
