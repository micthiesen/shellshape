package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"wine", "wine64"} {
		shellshape.Register(name, handleWine)
	}
}

// handleWine handles wine and wine64.
// The first positional is the program to run (an .exe path or Windows path) → <path>.
// Remaining positionals are arguments to the Windows program → classified with
// ClassifyToken, with unrecognized tokens collapsed to <val>.
// Boolean flags: --version.
func handleWine(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	var result []string
	exeFound := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			exeFound = true
			i++
			continue
		}

		if !exeFound && shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if !exeFound {
			// First positional is the exe/program path
			result = append(result, "<path>")
			exeFound = true
			i++
			continue
		}

		// Remaining args to the Windows program
		result = append(result, classifyWineArg(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// classifyWineArg classifies arguments passed to the Windows program.
// Windows-style paths (C:\...) and installer switches (/S, /D=...) become <val>.
// Unrecognized plain tokens also become <val> since they're program-specific data.
func classifyWineArg(tok string) string {
	// Windows installer switches like /S, /D=... → <val>
	if strings.HasPrefix(tok, "/") && !strings.Contains(tok[1:], "/") {
		return "<val>"
	}
	// Windows-style backslash paths
	if len(tok) >= 3 && tok[1] == ':' && tok[2] == '\\' {
		return "<val>"
	}
	classified := shellshape.ClassifyToken(tok)
	if classified != tok {
		return classified
	}
	return "<val>"
}
