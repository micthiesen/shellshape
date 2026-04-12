package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"zstd", "unzstd", "zstdcat", "zstdmt"} {
		shellshape.Register(name, handleZstd)
	}
}

func handleZstd(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	categories := []shellshape.FlagCategory{
		{Flags: map[string]bool{"-o": true, "-D": true}, Placeholder: "<path>"},
		{Flags: map[string]bool{"--threads": true}, Placeholder: "N"},
	}

	var result []string
	pathCount := 0
	firstPathIdx := -1

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// --fast=N fused flag
		if strings.HasPrefix(tok, "--fast=") {
			result = append(result, "--fast=N")
			i++
			continue
		}

		// Flag categories (consume next token)
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Compression level flags: -1 through -19
		if len(tok) >= 2 && tok[0] == '-' && tok[1] >= '1' && tok[1] <= '9' {
			allDigits := true
			for _, c := range tok[1:] {
				if c < '0' || c > '9' {
					allDigits = false
					break
				}
			}
			if allDigits {
				result = append(result, tok)
				i++
				continue
			}
		}

		// Other flags (boolean)
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: file path
		pathCount++
		if pathCount == 1 {
			firstPathIdx = len(result)
			result = append(result, shellshape.ClassifyToken(tok))
		} else if pathCount == 2 {
			result[firstPathIdx] = "<path>+"
		}
		// 3+ paths: already collapsed
		i++
	}

	result = append(result, redirects...)
	return result
}
