package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	shellshape.Register("mysql", handleMysql)
}

// handleMysql handles the mysql command.
// The first positional becomes <database>.
// -u/--user -> <user>, -h/--host -> <host>, -P/--port -> N,
// -e/--execute -> <query>, -S/--socket and other path flags -> <path>,
// -D/--database -> <database>.
// Bare -p is boolean (password prompt); fused -pVALUE collapses to -p<str>.
func handleMysql(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	userFlags := map[string]bool{
		"-u": true, "--user": true,
	}
	hostFlags := map[string]bool{
		"-h": true, "--host": true,
	}
	portFlags := map[string]bool{
		"-P": true, "--port": true,
	}
	queryFlags := map[string]bool{
		"-e": true, "--execute": true, "--init-command": true,
	}
	pathFlags := map[string]bool{
		"-S": true, "--socket": true,
		"--defaults-file": true, "--defaults-extra-file": true,
		"--ssl-ca": true, "--ssl-cert": true, "--ssl-key": true,
		"--tee": true, "--pager": true,
	}
	dbFlags := map[string]bool{
		"-D": true, "--database": true,
	}
	strFlags := map[string]bool{
		"--default-character-set": true, "--delimiter": true,
		"--default-auth": true, "--plugin-dir": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: userFlags, Placeholder: "<user>"},
		{Flags: hostFlags, Placeholder: "<host>"},
		{Flags: portFlags, Placeholder: "N"},
		{Flags: queryFlags, Placeholder: "<query>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: dbFlags, Placeholder: "<database>"},
		{Flags: strFlags, Placeholder: "<str>"},
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

		// Fused -pVALUE (password embedded in flag).
		if len(tok) > 2 && strings.HasPrefix(tok, "-p") && !strings.HasPrefix(tok, "--") {
			result = append(result, "-p<str>")
			i++
			continue
		}

		// Fused --flag=value
		if fused, ok := shellshape.ConsumeFusedFlag(tok, categories); ok {
			result = append(result, fused)
			i++
			continue
		}

		// Flag with separate next-token argument.
		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		// Boolean/standalone flags (including bare -p).
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: database name.
		result = append(result, "<database>")
		i++
	}

	result = append(result, redirects...)
	return result
}
