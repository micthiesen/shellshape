package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	for _, name := range []string{"sqlite3", "sqlite"} {
		shellshape.Register(name, handleSqlite3)
	}
}

// handleSqlite3 handles sqlite3 (and aliases like sqlite).
// Usage: sqlite3 [options] [database] [SQL]
// First positional is the database path (<path>).
// Second positional is the SQL query (<sql>).
// Flags that take a value argument: -separator, -newline, -cmd, -init.
func handleSqlite3(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valueFlagsSqlite := map[string]bool{
		"-separator": true, "-newline": true,
		"-cmd": true, "-init": true,
	}

	var result []string
	dbAssigned := false
	sqlAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !dbAssigned {
				dbAssigned = true
			} else {
				sqlAssigned = true
			}
			i++
			continue
		}

		if valueFlagsSqlite[tok] && i+1 < len(args) {
			result = append(result, tok, "<val>")
			i += 2
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positionals: first is database path, second is SQL
		if !dbAssigned {
			result = append(result, "<path>")
			dbAssigned = true
		} else if !sqlAssigned {
			result = append(result, "<sql>")
			sqlAssigned = true
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
