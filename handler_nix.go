package shellshape

import (
	"strings"
)

func init() {
	Register("nix", handleNix, HandlerOptions{HasSubcommands: true})
}

// handleNix handles the nix CLI (new-style nix3 commands).
// The first subcommand (build, shell, run, develop, flake, profile, etc.)
// is already consumed by the normalizer.
//
// Many subcommands take a second sub-subcommand (flake init, profile install,
// store gc) which is kept verbatim.
//
// Installable references (nixpkgs#hello, .#default, github:user/repo#pkg)
// are the primary data type and collapse to <pkg>.
//
// After "--" all tokens pass through verbatim (arguments to the program being run).
func handleNix(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"-o": true, "--out-link": true,
		"-I": true, "--include": true,
		"--output-lock-file":    true,
		"--reference-lock-file": true,
		"--profile":             true,
	}

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"--log-format":     true,
		"--inputs-from":    true,
		"--update-input":   true,
		"--eval-store":     true,
		"--to":             true,
		"--from":           true,
		"--command":        true,
		"--arg-from-stdin": true,
	}

	// Flags that consume the next TWO tokens as values.
	doubleArgFlags := map[string]bool{
		"--override-input": true,
		"--arg":            true,
		"--argstr":         true,
		"--arg-from-file":  true,
		"--option":         true,
	}

	// Expression flags: next token is a nix expression.
	exprFlags := map[string]bool{
		"--expr": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{valFlags, "<val>"},
		{exprFlags, "<expr>"},
	}

	// Subcommands that have their own sub-subcommand as the first positional.
	hasSubSub := map[string]bool{
		"flake":      true,
		"profile":    true,
		"store":      true,
		"registry":   true,
		"nar":        true,
		"key":        true,
		"derivation": true,
		"config":     true,
	}

	// Subcommands where all positionals after sub-subcommand are search
	// terms (collapse to <str> instead of <pkg>).
	searchLike := map[string]bool{
		"search": true,
	}

	var result []string
	subSubTaken := !hasSubSub[subcommand]
	firstSearchPositional := true // for search: first positional is the flake ref

	i := 0
	for i < len(args) {
		tok := args[i]

		// After "--", everything passes through verbatim.
		if tok == "--" {
			result = append(result, args[i:]...)
			break
		}

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Double-arg flags (consume TWO following tokens).
		if doubleArgFlags[tok] {
			result = append(result, tok)
			i++
			for n := 0; n < 2 && i < len(args); n++ {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// Single-arg flags.
		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Fused --flag=value.
		if fused, ok := consumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Boolean/unknown flags pass through.
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Sub-subcommand: first positional for subcommands that have them.
		if !subSubTaken {
			result = append(result, tok)
			subSubTaken = true
			i++
			continue
		}

		// Positional arguments.
		if searchLike[subcommand] {
			// For nix search: first positional is the flake source, rest are search terms.
			if firstSearchPositional {
				result = append(result, "<pkg>")
				firstSearchPositional = false
			} else {
				result = append(result, "<str>")
			}
		} else {
			result = append(result, nixClassifyPositional(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// nixClassifyPositional classifies a nix positional argument.
// Installable references (containing # or using flake URI schemes) become <pkg>.
// Everything else falls through to classifyToken.
func nixClassifyPositional(tok string) string {
	if isNixInstallable(tok) {
		return "<pkg>"
	}
	return classifyToken(tok)
}

// isNixInstallable returns true if the token looks like a nix installable reference.
// Patterns: "nixpkgs#hello", ".#default", "github:user/repo#pkg",
// "github:user/repo", "path:./foo", ".#packages.x86_64-linux.default"
func isNixInstallable(tok string) bool {
	// Contains "#" -> flake reference with attribute
	if strings.Contains(tok, "#") {
		return true
	}
	// Nix-specific URI schemes without #
	for _, prefix := range []string{"github:", "gitlab:", "sourcehut:", "path:", "git+", "tarball:", "file+"} {
		if strings.HasPrefix(tok, prefix) {
			return true
		}
	}
	return false
}
