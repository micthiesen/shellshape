package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("terraform", handleTerraform, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleTerraform handles terraform subcommand arguments.
// Since terraform is in subcommandExecutables, the subcommand (init, plan,
// apply, destroy, import, state, output, etc.) is already consumed before
// this handler is called.
//
// Terraform uses single-dash long flags. Value flags (-var, -target,
// -backend-config) collapse their argument to <val>. Path flags (-var-file,
// -state, -out) collapse to <path>. Numeric flags (-parallelism) collapse
// to N. Boolean flags (-auto-approve, -no-color, etc.) are kept verbatim.
// Flags can use either -flag value or -flag=value syntax.
// Remaining positionals use classifyToken, except for import where
// positionals (resource address, ID) collapse to <val>.
func handleTerraform(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Flags whose next token is a generic value.
	valFlags := map[string]bool{
		"-var":            true,
		"-target":         true,
		"-backend-config": true,
		"-plugin-dir":     true,
		"-lock-timeout":   true,
	}

	// Flags whose next token is a path.
	pathFlags := map[string]bool{
		"-var-file":  true,
		"-state":     true,
		"-state-out": true,
		"-out":       true,
		"-backup":    true,
		"-chdir":     true,
	}

	// Flags that are boolean (no argument consumed).
	boolFlags := map[string]bool{
		"-auto-approve":      true,
		"-no-color":          true,
		"-json":              true,
		"-compact-warnings":  true,
		"-detailed-exitcode": true,
		"-destroy":           true,
		"-reconfigure":       true,
		"-migrate-state":     true,
		"-upgrade":           true,
		"-force-copy":        true,
		"-get":               true,
		"-update":            true,
	}

	// Flags whose next token is numeric.
	numericFlags := map[string]bool{
		"-parallelism": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
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

		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Check for -flag=value syntax (single-dash terraform flags).
		if norm, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, norm)
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

		// Positional: import takes resource address + ID as values,
		// others use generic classification.
		if subcommand == "import" {
			result = append(result, "<val>")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
