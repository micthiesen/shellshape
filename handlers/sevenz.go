package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
	"strings"
)

func init() {
	for _, name := range []string{"7z", "7za", "7zr"} {
		shellshape.Register(name, handleSevenZ, shellshape.HandlerOptions{HasSubcommands: true})
	}
}

var mxLevelRE = regexp.MustCompile(`^-mx=\d+$`)

// handleSevenZ handles 7z, 7za, 7zr archive commands.
// Fused flags: -o<dir> → -o<path>, -p<pass> → -p<val>, -mx=N → -mx=N,
// -t<type> (structural, kept verbatim: -t7z, -tzip, -ttar).
// All positionals are file paths collapsed to <path>.
func handleSevenZ(subcommand string, tokens []string) []string {
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

		// Fused output dir: -o<path>
		if strings.HasPrefix(tok, "-o") && len(tok) > 2 {
			result = append(result, "-o<path>")
			i++
			continue
		}

		// Fused password: -p<val>
		if strings.HasPrefix(tok, "-p") && len(tok) > 2 {
			result = append(result, "-p<val>")
			i++
			continue
		}

		// Compression level: -mx=N
		if mxLevelRE.MatchString(tok) {
			result = append(result, "-mx=N")
			i++
			continue
		}

		// Archive type: -t<type> (structural, kept verbatim)
		if strings.HasPrefix(tok, "-t") && len(tok) > 2 {
			result = append(result, tok)
			i++
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: always a path (archive or files)
		result = append(result, "<path>")
		i++
	}

	result = append(result, redirects...)
	return result
}
