package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("psql", handlePsql)
}

// handlePsql handles the psql (PostgreSQL client) command.
// Connection Flags: -h (host), -p (port), -U (username), -d (dbname).
// Execution Flags: -c (query), -f (file path), -o/-L (output/log file paths).
// Formatting Flags: -F, -R, -P, -T, -v (opaque string values).
// First positional becomes <dbname>, second becomes <user>.
func handlePsql(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	hostFlags := map[string]bool{
		"-h": true, "--host": true,
	}
	portFlags := map[string]bool{
		"-p": true, "--port": true,
	}
	userFlags := map[string]bool{
		"-U": true, "--username": true,
	}
	dbnameFlags := map[string]bool{
		"-d": true, "--dbname": true,
	}
	queryFlags := map[string]bool{
		"-c": true, "--command": true,
	}
	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"-o": true, "--output": true,
		"-L": true, "--log-file": true,
	}
	strFlags := map[string]bool{
		"-F": true, "--field-separator": true,
		"-R": true, "--record-separator": true,
		"-P": true, "--pset": true,
		"-T": true, "--table-attr": true,
		"-v": true, "--set": true, "--variable": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: hostFlags, Placeholder: "<host>"},
		{Flags: portFlags, Placeholder: "N"},
		{Flags: userFlags, Placeholder: "<user>"},
		{Flags: dbnameFlags, Placeholder: "<dbname>"},
		{Flags: queryFlags, Placeholder: "<query>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: strFlags, Placeholder: "<str>"},
	}

	var result []string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		// --flag=value syntax
		if strings.HasPrefix(tok, "--") && strings.ContainsRune(tok, '=') {
			if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
				result = append(result, fused)
				i++
				continue
			}
			// Unknown long flag with =, keep as-is with generic classification
			result = append(result, shellshape.ClassifyToken(tok))
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

		// Positionals: first is dbname, second is username
		switch positionalCount {
		case 0:
			result = append(result, "<dbname>")
		default:
			result = append(result, "<user>")
		}
		positionalCount++
		i++
	}

	result = append(result, redirects...)
	return result
}
