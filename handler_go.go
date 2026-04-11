package shellshape

import "strings"

func init() {
	Register("go", handleGo, HandlerOptions{HasSubcommands: true})
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

func consumeFlagArg(tok string, args []string, i int, result []string, placeholder string) ([]string, int) {
	result = append(result, tok)
	i++
	if i < len(args) {
		if isSubshellToken(args[i]) {
			result = append(result, args[i])
		} else {
			result = append(result, placeholder)
		}
		i++
	}
	return result, i
}

// handleFusedFlag handles flags like -count=1, -timeout=5m where the value
// is joined with = sign. Returns the normalized token and true if it was a
// fused flag, or empty string and false otherwise.
func handleGoFusedFlag(tok string, numFlags, durationFlags, patternFlags, pathFlags, valFlags map[string]bool) (string, bool) {
	eqIdx := strings.IndexByte(tok, '=')
	if eqIdx < 0 {
		return "", false
	}
	key := tok[:eqIdx]
	if numFlags != nil && numFlags[key] {
		return key + "=N", true
	}
	if durationFlags != nil && durationFlags[key] {
		return key + "=<duration>", true
	}
	if patternFlags != nil && patternFlags[key] {
		return key + "=<pattern>", true
	}
	if pathFlags != nil && pathFlags[key] {
		return key + "=<path>", true
	}
	if valFlags != nil && valFlags[key] {
		return key + "=<val>", true
	}
	return "", false
}

func handleGoTest(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

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
		"-coverprofile":  true,
		"-cpuprofile":    true,
		"-memprofile":    true,
		"-mutexprofile":  true,
		"-blockprofile":  true,
		"-trace":         true,
		"-outputdir":     true,
		"-o":             true,
	}
	valFlags := map[string]bool{
		"-covermode":       true,
		"-coverpkg":        true,
		"-benchtime":       true,
		"-cpu":             true,
		"-blockprofilerate": true,
		"-memprofilerate":  true,
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

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Fused =value flags
		if strings.Contains(tok, "=") && isFlagToken(tok) {
			if norm, ok := handleGoFusedFlag(tok, numFlags, durationFlags, patternFlags, pathFlags, allValFlags); ok {
				result = append(result, norm)
				i++
				continue
			}
			if goBuildNumFlags[tok[:strings.IndexByte(tok, '=')]] {
				result = append(result, tok[:strings.IndexByte(tok, '=')]+"=N")
				i++
				continue
			}
			// Unknown fused flag, classify the value part
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		if patternFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<pattern>")
			continue
		}
		if numFlags[tok] || goBuildNumFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}
		if durationFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<duration>")
			continue
		}
		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}
		if allValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: package path
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoBuild(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true,
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

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}
		if goBuildValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoRun(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true,
	}

	var result []string
	i := 0
	seenPositional := false
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
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
				result, i = consumeFlagArg(tok, args, i, result, "<path>")
				continue
			}
			if goBuildValFlags[tok] {
				result, i = consumeFlagArg(tok, args, i, result, "<val>")
				continue
			}
			if goBuildNumFlags[tok] {
				result, i = consumeFlagArg(tok, args, i, result, "N")
				continue
			}
			if isFlagToken(tok) {
				result = append(result, tok)
				i++
				continue
			}
		}

		if !seenPositional {
			// First positional is the file/package to run
			seenPositional = true
			result = append(result, classifyToken(tok))
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
	args, redirects := splitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if goBuildValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Module paths like golang.org/x/tools/gopls@latest are structural
		// Relative paths like ./cmd/app classify normally
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, classifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoGet(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals are module paths — keep verbatim (structural)
		// unless they look like relative paths
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, classifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoMod(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// go mod edit flags that consume a value
	modEditValFlags := map[string]bool{
		"-replace":    true,
		"-require":    true,
		"-dropreplace": true,
		"-droprequire": true,
		"-exclude":    true,
		"-dropexclude": true,
		"-retract":    true,
		"-dropretract": true,
		"-go":         true,
		"-toolchain":  true,
		"-json":       true,
		"-fmt":        true,
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

		if modEditValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
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
	args, redirects := splitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if isSubshellToken(tok) {
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
	args, redirects := splitRedirects(tokens)

	var result []string
	for _, tok := range args {
		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		// -w flag with KEY=VALUE: collapse the value part
		if strings.Contains(tok, "=") && !isFlagToken(tok) {
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
	args, redirects := splitRedirects(tokens)

	valFlags := map[string]bool{
		"-f": true,
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

		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// "all" and module paths are structural; relative paths classify
		if strings.HasPrefix(tok, "./") || strings.HasPrefix(tok, "../") {
			result = append(result, classifyToken(tok))
		} else {
			result = append(result, tok)
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

func handleGoTool(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	seenTool := false
	for i := 0; i < len(args); i++ {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		if !seenTool {
			// First positional is the tool name (pprof, trace, dist) — structural
			result = append(result, tok)
			seenTool = true
			continue
		}

		// After tool name, classify remaining tokens
		result = append(result, classifyToken(tok))
	}

	result = append(result, redirects...)
	return result
}

func handleGoDefault(tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if goBuildValFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}
		if goBuildNumFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "N")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
