package shellshape

import "strings"

func init() {
	Register("dotnet", handleDotnet, HandlerOptions{HasSubcommands: true})
}

// dotnet flags that consume a path argument.
var dotnetPathFlags = map[string]bool{
	"-o": true, "--output": true,
	"--project":           true,
	"--artifacts-path":    true,
	"--results-directory": true,
	"--settings":          true,
}

// dotnet flags that consume a generic value argument.
var dotnetValFlags = map[string]bool{
	"-c": true, "--configuration": true,
	"-f": true, "--framework": true,
	"-r": true, "--runtime": true,
	"-v": true, "--verbosity": true,
	"--arch": true, "--os": true,
	"--source": true, "-s": true,
	"--version":        true,
	"--version-suffix": true,
	"-n":               true, "--name": true,
	"--logger":             true,
	"--collect":            true,
	"--blame-hang-timeout": true,
	"--launch-profile":     true,
	"--api-key":            true,
	"--roll-forward":       true,
	"--fx-version":         true,
}

// dotnet flags that consume a pattern argument.
var dotnetPatternFlags = map[string]bool{
	"--filter": true,
}

var dotnetCategories = []flagCategory{
	{dotnetPatternFlags, "<pattern>"},
	{dotnetPathFlags, "<path>"},
	{dotnetValFlags, "<val>"},
}

// handleDotnet normalizes dotnet CLI invocations.
// The subcommand (build, run, test, etc.) is already extracted by the framework.
func handleDotnet(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// After "--", all remaining tokens are pass-through application arguments.
	var passthrough []string
	for j, tok := range args {
		if tok == "--" {
			passthrough = args[j:]
			args = args[:j]
			break
		}
	}

	// Track whether we've seen "package" or "reference" as a sub-subcommand
	// in "add"/"remove" subcommands, to keep package names structural.
	seenSubSub := false
	keepNextVerbatim := false

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Handle -p:Key=Value and --property:Key=Value (MSBuild properties).
		if strings.HasPrefix(tok, "-p:") || strings.HasPrefix(tok, "--property:") {
			colonIdx := strings.IndexByte(tok, ':')
			result = append(result, tok[:colonIdx+1]+"<val>")
			i++
			continue
		}

		// Fused --flag=value syntax.
		if strings.Contains(tok, "=") && isFlagToken(tok) {
			if norm, ok := consumeFusedFlag(tok, dotnetCategories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, dotnetCategories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// For "add"/"remove" subcommands, detect "package" or "reference" sub-subcommand.
		if !seenSubSub && (subcommand == "add" || subcommand == "remove") {
			if tok == "package" {
				seenSubSub = true
				keepNextVerbatim = true // package name is structural
				result = append(result, tok)
				i++
				continue
			}
			if tok == "reference" {
				seenSubSub = true
				result = append(result, tok)
				i++
				continue
			}
		}

		// Keep package/tool names verbatim.
		if keepNextVerbatim {
			keepNextVerbatim = false
			result = append(result, tok)
			i++
			continue
		}

		// Positional arguments: classify based on context.
		result = append(result, dotnetClassifyPositional(subcommand, tok))
		i++
	}

	// Append passthrough args verbatim (including the "--" separator).
	result = append(result, passthrough...)

	result = append(result, redirects...)
	return result
}

// dotnetClassifyPositional determines how to classify a positional argument.
// Package names, tool names, and template names are structural (kept verbatim).
// File paths and other data are collapsed.
func dotnetClassifyPositional(subcommand string, tok string) string {
	// .dll files are always paths (not in the default pathExtensionsRE).
	if strings.HasSuffix(tok, ".dll") {
		return "<path>"
	}

	// .nupkg files are paths.
	if strings.HasSuffix(tok, ".nupkg") {
		return "<path>"
	}

	return classifyToken(tok)
}
