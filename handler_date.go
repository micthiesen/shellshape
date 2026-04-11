package shellshape

import "strings"

func init() {
	Register("date", handleDate)
}

// handleDate handles the date command.
// Format strings (positionals starting with +) become <fmt>.
// -d/--date and -s/--set consume the next arg as <date-str>.
// -r consumes the next arg as <date-str> (macOS: timestamp, Linux: ref file).
// -f/--file consumes the next arg as <path>.
// --date=VAL and --set=VAL become --date=<date-str> and --set=<date-str>.
// --rfc-3339=TIMESPEC is kept verbatim (the value is a keyword).
// Other positionals become <date-str>.
func handleDate(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	dateStrFlags := map[string]bool{
		"-d": true, "--date": true,
		"-s": true, "--set": true,
		"-r": true, "--reference": true,
	}

	fileFlags := map[string]bool{
		"-f": true, "--file": true,
	}

	// Long flags with = that get <date-str> replacement
	dateStrLongFlags := map[string]bool{
		"--date": true, "--set": true, "--reference": true,
	}

	// Long flags with = that keep their value verbatim (keyword args)
	verbatimLongFlags := map[string]bool{
		"--rfc-3339": true, "--iso-8601": true,
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

		// Handle long flags with =
		if strings.HasPrefix(tok, "--") {
			if idx := strings.Index(tok, "="); idx >= 0 {
				name := tok[:idx]
				if dateStrLongFlags[name] {
					result = append(result, name+"=<date-str>")
					i++
					continue
				}
				if verbatimLongFlags[name] {
					result = append(result, tok)
					i++
					continue
				}
				// Other long=val flags: collapse value
				result = append(result, classifyToken(tok))
				i++
				continue
			}
		}

		// Short/long flags that consume next token as date string
		if dateStrFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<date-str>")
				i++
			}
			continue
		}

		// Flags that consume next token as path
		if fileFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: format string starts with +
		if strings.HasPrefix(tok, "+") {
			result = append(result, "<fmt>")
			i++
			continue
		}

		// Other positional: date string for setting the clock
		result = append(result, "<date-str>")
		i++
	}

	result = append(result, redirects...)
	return result
}
