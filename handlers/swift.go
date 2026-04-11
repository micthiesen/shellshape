package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("swift", handleSwift)
	shellshape.Register("swiftc", handleSwiftc)
}

// Known swift subcommands.
var swiftSubcommands = map[string]bool{
	"build": true, "test": true, "run": true,
	"package": true, "repl": true,
}

// Compiler flags shared by swift (no-subcommand mode) and swiftc.
var swiftCompilerPathFlags = map[string]bool{
	"-o": true, "-sdk": true,
	"-module-cache-path": true, "-working-directory": true,
	"-cas-path": true, "-cas-plugin-path": true,
	"-access-notes-path": true, "-file-compilation-dir": true,
	"-dependency-scan-serialize-diagnostics-path": true,
	"-I": true, "-L": true, "-F": true, "-Fsystem": true,
}

var swiftCompilerValFlags = map[string]bool{
	"-target": true, "-module-name": true, "-module-link-name": true,
	"-swift-version": true, "-assert-config": true,
	"-D": true, "-e": true,
	"-Xcc": true, "-Xlinker": true, "-Xswiftc": true,
	"-framework": true, "-l": true,
	"-clang-target": true, "-clang-target-variant": true,
	"-diagnostic-style": true, "-export-as": true,
	"-embed-tbd-for-module":        true,
	"-enable-experimental-feature": true, "-disable-experimental-feature": true,
	"-enable-upcoming-feature": true, "-disable-upcoming-feature": true,
	"-allowable-client":          true,
	"-default-isolation":         true,
	"-vfsoverlay":                true,
	"-cxx-interoperability-mode": true,
	"-debug-info-format":         true,
	"-enforce-exclusivity":       true,
	"-dwarf-version":             true,
}

var swiftCompilerNumFlags = map[string]bool{
	"-num-threads": true,
}

// SPM flags shared across swift build/test/run/package subcommands.
var swiftSPMPathFlags = map[string]bool{
	"--package-path":    true,
	"--cache-path":      true,
	"--config-path":     true,
	"--security-path":   true,
	"--scratch-path":    true,
	"--swift-sdks-path": true,
	"--netrc-file":      true,
}

var swiftSPMValFlags = map[string]bool{
	"-c": true, "--configuration": true,
	"--swift-sdk":                        true,
	"--traits":                           true,
	"--manifest-cache":                   true,
	"--toolset":                          true,
	"--pkg-config-path":                  true,
	"--resolver-fingerprint-checking":    true,
	"--resolver-signing-entity-checking": true,
}

var swiftSPMNumFlags = map[string]bool{
	"-j": true, "--jobs": true,
}

func handleSwift(_ string, tokens []string) []string {
	// Manually extract subcommand since swift can also act as a compiler.
	if len(tokens) > 0 && swiftSubcommands[tokens[0]] {
		sub := tokens[0]
		rest := tokens[1:]
		switch sub {
		case "build":
			return append([]string{sub}, handleSwiftBuild(rest)...)
		case "test":
			return append([]string{sub}, handleSwiftTest(rest)...)
		case "run":
			return append([]string{sub}, handleSwiftRun(rest)...)
		case "package":
			return append([]string{sub}, handleSwiftPackage(rest)...)
		case "repl":
			return append([]string{sub}, handleSwiftRepl(rest)...)
		}
	}
	// No subcommand: swift acts as a compiler (like swiftc)
	return handleSwiftCompiler(tokens)
}

func handleSwiftc(subcommand string, tokens []string) []string {
	return handleSwiftCompiler(tokens)
}

func handleSwiftCompiler(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: swiftCompilerPathFlags, Placeholder: "<path>"},
		{Flags: swiftCompilerValFlags, Placeholder: "<val>"},
		{Flags: swiftCompilerNumFlags, Placeholder: "N"},
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

		if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, shellshape.ClassifyToken(tok))
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

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleSwiftBuild(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	buildValFlags := map[string]bool{
		"--product": true,
		"--target":  true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: swiftSPMPathFlags, Placeholder: "<path>"},
		{Flags: swiftSPMValFlags, Placeholder: "<val>"},
		{Flags: buildValFlags, Placeholder: "<val>"},
		{Flags: swiftSPMNumFlags, Placeholder: "N"},
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

		if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, shellshape.ClassifyToken(tok))
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

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleSwiftTest(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	testPathFlags := map[string]bool{
		"--xunit-output":     true,
		"--attachments-path": true,
	}
	testPatternFlags := map[string]bool{
		"--filter": true,
		"--skip":   true,
	}
	testValFlags := map[string]bool{
		"-s": true, "--specifier": true,
		"--product": true,
		"--target":  true,
	}
	testNumFlags := map[string]bool{
		"--num-workers": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: swiftSPMPathFlags, Placeholder: "<path>"},
		{Flags: testPathFlags, Placeholder: "<path>"},
		{Flags: testPatternFlags, Placeholder: "<pattern>"},
		{Flags: swiftSPMValFlags, Placeholder: "<val>"},
		{Flags: testValFlags, Placeholder: "<val>"},
		{Flags: swiftSPMNumFlags, Placeholder: "N"},
		{Flags: testNumFlags, Placeholder: "N"},
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

		if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, shellshape.ClassifyToken(tok))
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

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleSwiftRun(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	runValFlags := map[string]bool{
		"--product": true,
		"--target":  true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: swiftSPMPathFlags, Placeholder: "<path>"},
		{Flags: swiftSPMValFlags, Placeholder: "<val>"},
		{Flags: runValFlags, Placeholder: "<val>"},
		{Flags: swiftSPMNumFlags, Placeholder: "N"},
	}

	var result []string
	i := 0
	seenExecutable := false
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			if !seenExecutable {
				seenExecutable = true
			}
			continue
		}

		// Before the executable name, parse SPM flags
		if !seenExecutable {
			if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
				if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
					result = append(result, norm)
					i++
					continue
				}
				result = append(result, shellshape.ClassifyToken(tok))
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

			// First positional is the executable name — keep verbatim (structural)
			seenExecutable = true
			result = append(result, tok)
			i++
			continue
		}

		// After the executable name, all remaining tokens are program arguments — keep verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleSwiftPackage(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pkgValFlags := map[string]bool{
		"--type": true,
		"--name": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: swiftSPMPathFlags, Placeholder: "<path>"},
		{Flags: swiftSPMValFlags, Placeholder: "<val>"},
		{Flags: pkgValFlags, Placeholder: "<val>"},
		{Flags: swiftSPMNumFlags, Placeholder: "N"},
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

		if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, shellshape.ClassifyToken(tok))
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

		// Subcommand words (init, update, resolve, clean, etc.) are structural
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleSwiftRepl(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: swiftSPMPathFlags, Placeholder: "<path>"},
		{Flags: swiftSPMValFlags, Placeholder: "<val>"},
		{Flags: swiftSPMNumFlags, Placeholder: "N"},
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

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
