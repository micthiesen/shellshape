package shellshape

import "strings"

func init() {
	Register("java", handleJava)
	Register("javac", handleJavac)
}

// handleJava handles the java command.
// -jar makes the next arg a path; -cp/-classpath/--class-path take a path;
// -D, -X, -javaagent, -agentlib, -agentpath are fused prefix flags;
// --add-modules/opens/reads/exports and -m/--module take a value.
// First non-flag positional is the class name (kept verbatim); subsequent
// positionals are program arguments collapsed to <arg>.
func handleJava(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-cp": true, "-classpath": true, "--class-path": true,
		"--module-path": true, "-p": true,
		"--upgrade-module-path": true,
		"--patch-module":        true,
	}

	valFlags := map[string]bool{
		"--add-modules": true, "--add-opens": true,
		"--add-reads": true, "--add-exports": true,
		"-m": true, "--module": true,
	}

	// -jar consumes the next token as a path
	jarFlag := "-jar"

	var result []string
	classSeen := false // true after class name or -jar <path>

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			if !classSeen {
				classSeen = true
			}
			i++
			continue
		}

		// After class name or -jar <path>, everything is a program arg
		if classSeen {
			result = append(result, "<arg>")
			i++
			continue
		}

		// Fused prefix flags: -Dkey=val, -Dfoo, -Xms256m, -Xmx2g,
		// -javaagent:path, -agentlib:opts, -agentpath:path
		if strings.HasPrefix(tok, "-D") && len(tok) > 2 {
			result = append(result, "-D<val>")
			i++
			continue
		}
		if strings.HasPrefix(tok, "-X") && len(tok) > 2 {
			result = append(result, "-X<val>")
			i++
			continue
		}
		if strings.HasPrefix(tok, "-javaagent:") {
			result = append(result, "-javaagent:<val>")
			i++
			continue
		}
		if strings.HasPrefix(tok, "-agentlib:") {
			result = append(result, "-agentlib:<val>")
			i++
			continue
		}
		if strings.HasPrefix(tok, "-agentpath:") {
			result = append(result, "-agentpath:<val>")
			i++
			continue
		}

		// -jar flag: next token is a jar path
		if tok == jarFlag {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			classSeen = true
			continue
		}

		if pathFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valFlags[tok] {
			result, i = consumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// First positional: class name, keep verbatim (structural)
		result = append(result, tok)
		classSeen = true
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleJavac handles the javac compiler.
// Path-consuming flags: -cp, -classpath, --class-path, -sourcepath, --source-path,
// -d, -s, -h, --module-path, -p, -processorpath, --processor-path,
// --processor-module-path, -bootclasspath, --system, --upgrade-module-path.
// Value-consuming flags: -source, --source, -target, --target, --release,
// -encoding, -processor, --module, -m, --add-modules, --add-opens, etc.
// All positionals are source files, collapsed via classifyToken.
func handleJavac(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-cp": true, "-classpath": true, "--class-path": true,
		"-sourcepath": true, "--source-path": true,
		"-d": true, "-s": true, "-h": true,
		"--module-path": true, "-p": true,
		"-processorpath": true, "--processor-path": true,
		"--processor-module-path": true,
		"-bootclasspath":          true,
		"--system":                true,
		"--upgrade-module-path":   true,
		"-endorseddirs":           true,
		"-extdirs":                true,
	}

	valFlags := map[string]bool{
		"-source": true, "--source": true,
		"-target": true, "--target": true,
		"--release":     true,
		"-encoding":     true,
		"-processor":    true,
		"--module":      true,
		"-m":            true,
		"--add-modules": true,
		"--add-opens":   true,
		"--add-reads":   true,
		"--add-exports": true,
		"-profile":      true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{valFlags, "<val>"},
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
			if norm, ok := consumeFusedFlag(tok, categories); ok {
				result = append(result, norm)
				i++
				continue
			}
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: source file
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
