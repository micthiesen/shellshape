package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"tsc", "tsgo"} {
		shellshape.Register(name, handleTsc)
	}
}

// handleTsc handles tsc (TypeScript compiler).
// Path flags (-p, --project, --outDir, etc.) consume the next token as <path>.
// Value flags (--target, --module, etc.) consume the next token as <val>.
// Remaining positionals are file paths classified via classifyToken.
func handleTsc(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-p": true, "--project": true,
		"--outDir": true, "--rootDir": true, "--outFile": true,
		"--declarationDir": true, "--baseUrl": true,
		"--rootDirs": true, "--typeRoots": true,
		"--configFilePath": true,
	}

	valueFlags := map[string]bool{
		"--target": true, "--module": true, "--moduleResolution": true,
		"--lib": true, "--jsx": true, "--jsxFactory": true,
		"--jsxFragmentFactory": true, "--jsxImportSource": true,
		"--plugins": true, "--newLine": true,
		"--importsNotUsedAsValues": true, "--moduleSuffixes": true,
		"--moduleDetection": true, "--verbatimModuleSyntax": true,
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

		if pathFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<path>")
			continue
		}

		if valueFlags[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, shellshape.ClassifyToken(tok))
			i++
			continue
		}

		// Positional: file path
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
