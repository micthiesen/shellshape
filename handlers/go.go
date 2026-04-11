package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("go", handleGo, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleGo handles the go tool's subcommands (test, build, run, get, mod, etc.).
// The subcommand is already extracted by the framework.
//
// Key design: go test flags like -run, -bench collapse their pattern argument;
// build flags like -o, -ldflags collapse their value; go get/mod/doc/env keep
// module paths and symbol names verbatim since they're structural.
func handleGo(subcommand string, tokens []string) []string {
	switch subcommand {
	case "test":
		return handleGoTest(tokens)
	case "build":
		return handleGoBuild(tokens)
	case "install":
		return handleGoInstall(tokens)
	case "run":
		return handleGoRun(tokens)
	case "get":
		return handleGoGet(tokens)
	case "mod":
		return handleGoMod(tokens)
	case "doc":
		return handleGoDoc(tokens)
	case "env":
		return handleGoEnv(tokens)
	case "list":
		return handleGoList(tokens)
	case "tool":
		return handleGoTool(tokens)
	default:
		// fmt, vet, generate, clean, version, help, etc.
		return handleGoDefault(tokens)
	}
}

// Flags shared across build/test/run/vet that consume a value argument.
var goBuildValFlags = map[string]bool{
	"-ldflags":    true,
	"-gcflags":    true,
	"-asmflags":   true,
	"-gccgoflags": true,
	"-tags":       true,
	"-mod":        true,
	"-modfile":    true,
	"-overlay":    true,
	"-pgo":        true,
	"-pkgdir":     true,
	"-toolexec":   true,
}

// Flags shared across build/test/run/vet that consume a numeric argument.
var goBuildNumFlags = map[string]bool{
	"-p": true,
}

func handleGoTest(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	patternFlags := map[string]bool{
		"-run":   true,
		"-bench": true,
		"-fuzz":  true,
	}
	numFlags := map[string]bool{
		"-count":    true,
		"-parallel": true,
	}
	durationFlags := map[string]bool{
		"-timeout":  true,
		"-fuzztime": true,
	}
	pathFlags := map[string]bool{
		"-coverprofile": true,
		"-cpuprofile":   true,
		"-memprofile":   true,
		"-mutexprofile": true,
		"-blockprofile": true,
		"-trace":        true,
		"-outputdir":    true,
		"-o":            true,
	}
	valFlags := map[string]bool{
		"-covermode":            true,
		"-coverpkg":             true,
		"-benchtime":            true,
		"-cpu":                  true,
		"-blockprofilerate":     true,
		"-memprofilerate":       true,
		"-mutexprofilefraction": true,
	}
	// merge build val flags
	allValFlags := map[string]bool{}
	for k := range valFlags {
		allValFlags[k] = true
	}
	for k := range goBuildValFlags {
		allValFlags[k] = true
	}

	categories := []shellshape.FlagCategory{
		{Flags: patternFlags, Placeholder: "<pattern>"},
		{Flags: numFlags, Placeholder: "N"},
		{Flags: goBuildNumFlags, Placeholder: "N"},
		{Flags: durationFlags, Placeholder: "<duration>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: allValFlags, Placeholder: "<val>"},
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

		// Fused =value flags
		if strings.Contains(tok, "=") && shellshape.IsFlagToken(tok) {
			if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			// Unknown fused flag, classify the value part
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

		// Positional: package path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoBuild(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true,
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

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}
		if goBuildValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
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

func handleGoRun(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true,
	}

	var result []string
	i := 0
	seenPositional := false
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			if !seenPositional {
				seenPositional = true
			}
			continue
		}

		// Before the first positional, parse go flags
		if !seenPositional {
			if pathFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
				continue
			}
			if goBuildValFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
				continue
			}
			if goBuildNumFlags[tok] {
				result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
				continue
			}
			if shellshape.IsFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}
		}

		if !seenPositional {
			// First positional is the file/package to run
			seenPositional = true
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// After the file/package, remaining tokens are program arguments
		// Keep them verbatim as they're part of the user's program invocation
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoInstall(tokens []string) []string {
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

		if goBuildValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Module paths like golang.org/x/tools/gopls@latest are structural
		// Relative paths like ./cmd/app classify normally
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, shellshape.ClassifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoGet(tokens []string) []string {
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

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals are module paths — keep verbatim (structural)
		// unless they look like relative paths
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, shellshape.ClassifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoMod(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// go mod edit flags that consume a value
	modEditValFlags := map[string]bool{
		"-replace":     true,
		"-require":     true,
		"-dropreplace": true,
		"-droprequire": true,
		"-exclude":     true,
		"-dropexclude": true,
		"-retract":     true,
		"-dropretract": true,
		"-go":          true,
		"-toolchain":   true,
		"-json":        true,
		"-fmt":         true,
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

		if modEditValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Module paths and subcommand words are structural
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoDoc(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}
		// All tokens are either boolean flags or symbol names (structural)
		result = append(result, tok)
	}

	result = append(result, redirects...)
	return result
}

func handleGoEnv(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// -w flag with KEY=VALUE: collapse the value part
		if strings.Contains(tok, "=") && !shellshape.IsFlagToken(tok) {
			eqIdx := strings.IndexByte(tok, '=')
			key := tok[:eqIdx]
			result = append(result, key+"=<val>")
			continue
		}

		// Env var names and flags are structural
		result = append(result, tok)
	}

	result = append(result, redirects...)
	return result
}

func handleGoList(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-f": true,
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

		if valFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// "all" and module paths are structural; relative paths classify
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, shellshape.ClassifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// Flags for go tool subcommands (cover, pprof, trace) that consume a path arg.
var goToolPathFlags = map[string]bool{
	"-func": true, "-html": true, "-o": true,
	"-http": true, // pprof: -http :8080 (actually a listen addr, but collapse)
}

func handleGoTool(tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	seenTool := false
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle fused -flag=value (e.g. -func=/tmp/cov.out)
		if shellshape.IsFlagToken(tok) && strings.Contains(tok, "=") {
			eqIdx := strings.IndexByte(tok, '=')
			key := tok[:eqIdx]
			if goToolPathFlags[key] {
				result = append(result, key+"=<path>")
			} else {
				result = append(result, shellshape.ClassifyToken(tok))
			}
			i++
			continue
		}

		// Handle separate -flag value
		if goToolPathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if !seenTool {
			// First positional is the tool name (cover, pprof, trace, dist) — structural
			result = append(result, tok)
			seenTool = true
			i++
			continue
		}

		// After tool name, classify remaining tokens
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoDefault(tokens []string) []string {
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

		if goBuildValFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "N")
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
