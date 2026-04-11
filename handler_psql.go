package shellshape

import "strings"

func init() {
	Register("psql", handlePsql)
}

// handlePsql handles the psql (PostgreSQL client) command.
// Connection flags: -h (host), -p (port), -U (username), -d (dbname).
// Execution flags: -c (query), -f (file path), -o/-L (output/log file paths).
// Formatting flags: -F, -R, -P, -T, -v (opaque string values).
// First positional becomes <dbname>, second becomes <user>.
func handlePsql(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

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

	categories := []flagCategory{
		{hostFlags, "<host>"},
		{portFlags, "N"},
		{userFlags, "<user>"},
		{dbnameFlags, "<dbname>"},
		{queryFlags, "<query>"},
		{pathFlags, "<path>"},
		{strFlags, "<str>"},
	}

	var result []string
	positionalCount := 0

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			positionalCount++
			i++
			continue
		}

		// --flag=value syntax
		if strings.HasPrefix(tok, "--") && strings.ContainsRune(tok, '=') {
			if fused, ok := consumeFusedFlag(tok, categories); ok {
				result = append(result, fused)
				i++
				continue
			}
			// Unknown long flag with =, keep as-is with generic classification
			result = append(result, classifyToken(tok))
			i++
			continue
		}

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
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
