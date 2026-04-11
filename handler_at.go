package shellshape

func init() {
	Register("at", handleAt)
	Register("batch", handleAt)
	Register("atq", handleAtq)
	Register("atrm", handleAtrm)
}

// handleAt handles the at and batch commands.
// Flags: -f (file), -q (queue), -t (time), -c (cat job ID).
// -d switches to delete mode where positionals are job IDs.
// Boolean flags: -m, -M, -v, -l, -b.
// All other positionals form the time specification and collapse to <time>.
func handleAt(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	fileFlags := map[string]bool{"-f": true}
	queueFlags := map[string]bool{"-q": true}
	timeArgFlags := map[string]bool{"-t": true}
	jobFlags := map[string]bool{"-c": true}

	deleteMode := false

	var result []string
	hasTime := false
	hasJobID := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if fileFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if queueFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<queue>")
				i++
			}
			continue
		}

		if timeArgFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<time>")
				i++
			}
			continue
		}

		if jobFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<job-id>")
				i++
			}
			continue
		}

		if tok == "-d" {
			result = append(result, tok)
			deleteMode = true
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		if deleteMode {
			if !hasJobID {
				result = append(result, "<job-id>")
				hasJobID = true
			}
		} else {
			if !hasTime {
				result = append(result, "<time>")
				hasTime = true
			}
		}
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleAtq handles the atq command (list pending jobs).
// Only -q (queue) consumes the next token; everything else is boolean.
func handleAtq(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if tok == "-q" {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<queue>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// atq has no meaningful positionals
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleAtrm handles the atrm command (remove scheduled jobs).
// All positionals are job IDs that collapse to a single <job-id>.
func handleAtrm(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	var result []string
	hasJobID := false

	for _, tok := range args {
		if isSubshellToken(tok) {
			result = append(result, tok)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			continue
		}

		// All positionals are job IDs
		if !hasJobID {
			result = append(result, "<job-id>")
			hasJobID = true
		}
	}

	result = append(result, redirects...)
	return result
}
