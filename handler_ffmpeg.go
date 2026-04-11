package shellshape

import "strings"

func init() {
	for _, name := range []string{"ffmpeg", "ffprobe"} {
		Register(name, handleFfmpeg)
	}
}

// handleFfmpeg handles ffmpeg and ffprobe commands.
// -i (input path) can appear multiple times; each consumes the next arg as <path>.
// Codec flags (-c, -c:v, -c:a, -codec, -vcodec, -acodec) keep their value verbatim
// because the codec choice is structural to the command's shape.
// -f (format) also keeps its value verbatim.
// -v (loglevel) keeps its value verbatim.
// Filter flags (-vf, -af, -filter, -filter:v, -filter:a) collapse to <filter>.
// Numeric flags (-r, -ar, -ac) collapse to N.
// Time/size flags (-ss, -to, -t, -s, -aspect, -b:v, -b:a, -ab, -b) collapse to <val>.
// -map collapses to <val>.
// -metadata collapses to <val>.
// Boolean flags (-y, -n, -vn, -an, -sn, -stats) take no argument.
// Remaining positionals are classified as paths.
func handleFfmpeg(_ string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags whose next arg is a path
	pathFlags := map[string]bool{
		"-i": true,
	}

	// Flags whose next arg is kept verbatim (structural choice)
	verbatimFlags := map[string]bool{
		"-f": true, "-v": true,
	}

	// Flags whose next arg collapses to <filter>
	filterFlags := map[string]bool{
		"-vf": true, "-af": true,
	}

	// Flags whose next arg collapses to N
	numericFlags := map[string]bool{
		"-r": true, "-ar": true, "-ac": true,
	}

	// Flags whose next arg collapses to <val>
	valFlags := map[string]bool{
		"-ss": true, "-to": true, "-t": true,
		"-s": true, "-aspect": true,
		"-map": true, "-metadata": true,
		"-ab": true, "-aq": true,
	}

	// Boolean flags (no argument consumed)
	boolFlags := map[string]bool{
		"-y": true, "-n": true,
		"-vn": true, "-an": true, "-sn": true,
		"-stats": true,
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
			result = append(result, tok)
			i++
			if i < len(args) {
				if isSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, "<path>")
				}
				i++
			}
			continue
		}

		// Codec flags: -c, -c:v, -c:a, -c:<anything>, -codec, -vcodec, -acodec
		if isCodecFlag(tok) {
			result = append(result, tok)
			i++
			if i < len(args) {
				// Keep codec value verbatim
				result = append(result, args[i])
				i++
			}
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Bitrate flags: -b, -b:v, -b:a
		if isBitrateFlag(tok) {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		// Filter flags: -filter, -filter:v, -filter:a, -filter:<anything>
		if filterFlags[tok] || isFilterFlag(tok) {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<filter>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if boolFlags[tok] {
			result = append(result, tok)
			i++
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: treat as path
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}

// isCodecFlag returns true for -c, -c:v, -c:a, -c:<spec>, -codec, -vcodec, -acodec.
func isCodecFlag(tok string) bool {
	if tok == "-c" || tok == "-codec" || tok == "-vcodec" || tok == "-acodec" {
		return true
	}
	if strings.HasPrefix(tok, "-c:") {
		return true
	}
	return false
}

// isBitrateFlag returns true for -b, -b:v, -b:a, -b:<spec>.
func isBitrateFlag(tok string) bool {
	if tok == "-b" {
		return true
	}
	if strings.HasPrefix(tok, "-b:") {
		return true
	}
	return false
}

// isFilterFlag returns true for -filter, -filter:v, -filter:a, -filter:<spec>.
func isFilterFlag(tok string) bool {
	if tok == "-filter" {
		return true
	}
	if strings.HasPrefix(tok, "-filter:") {
		return true
	}
	return false
}
