package shellshape

import "strings"

func init() {
	for _, name := range []string{"gradle", "gradlew"} {
		Register(name, handleGradle)
	}
}

// handleGradle handles gradle and gradlew.
// Task names (positionals) are structural and kept verbatim.
// -x/--exclude-task consumes the next token as <task>.
// Path flags (-b, -p, -c, -g, -I, etc.) collapse to <path>.
// Value flags (--type, --priority, --console, --warning-mode) collapse to <val>.
// Numeric flags (--max-workers) collapse to N.
// -D and -P prefixed tokens (system/project properties) collapse to <sysprop>/<projprop>.
func handleGradle(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	taskFlags := map[string]bool{
		"-x": true, "--exclude-task": true,
	}

	pathFlags := map[string]bool{
		"-b": true, "--build-file": true,
		"-c": true, "--settings-file": true,
		"-g": true, "--gradle-user-home": true,
		"-p": true, "--project-dir": true,
		"-I": true, "--init-script": true,
		"--include-build":     true,
		"--project-cache-dir": true,
	}

	valFlags := map[string]bool{
		"--type":                         true,
		"--priority":                     true,
		"--console":                      true,
		"--warning-mode":                 true,
		"--configuration-cache-problems": true,
	}

	numericFlags := map[string]bool{
		"--max-workers": true,
	}

	categories := []flagCategory{
		{taskFlags, "<task>"},
		{pathFlags, "<path>"},
		{valFlags, "<val>"},
		{numericFlags, "N"},
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

		// -Dproperty=value → <sysprop>
		if strings.HasPrefix(tok, "-D") && len(tok) > 2 {
			result = append(result, "<sysprop>")
			i++
			continue
		}

		// -Pproperty=value → <projprop>
		if strings.HasPrefix(tok, "-P") && len(tok) > 2 {
			result = append(result, "<projprop>")
			i++
			continue
		}

		// Fused --flag=value
		if fused, ok := consumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag with argument
		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Boolean or unknown flags kept verbatim
		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: task name, kept verbatim
		result = append(result, tok)
		i++
	}

	result = append(result, redirects...)
	return result
}
