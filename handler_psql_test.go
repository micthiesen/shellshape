package shellshape

import "testing"

func TestPsql(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"database only", "psql mydb", "psql <dbname>"},
		{"database and user", "psql mydb admin", "psql <dbname> <user>"},
		{"no args", "psql", "psql"},
		{"list databases", "psql -l", "psql -l"},

		// Connection flags
		{"host and port", "psql -h localhost -p 5432 mydb", "psql -h <host> -p N <dbname>"},
		{"full connection", "psql -h db.example.com -p 5432 -U admin -W mydb", "psql -h <host> -p N -U <user> -W <dbname>"},
		{"dbname flag", "psql -d postgres", "psql -d <dbname>"},
		{"username flag", "psql -U postgres mydb", "psql -U <user> <dbname>"},

		// Query and file execution
		{"command flag", "psql -c 'SELECT 1' mydb", "psql -c <query> <dbname>"},
		{"file flag", "psql mydb -f schema.sql", "psql <dbname> -f <path>"},
		{"output flag", "psql -o results.txt -c 'SELECT 1' mydb", "psql -o <path> -c <query> <dbname>"},
		{"log file", "psql -L session.log mydb", "psql -L <path> <dbname>"},

		// Formatting flags
		{"html output", "psql -H -c 'SELECT 1' mydb", "psql -H -c <query> <dbname>"},
		{"tuples only aligned", "psql -tA mydb", "psql -tA <dbname>"},
		{"field separator", "psql -F ',' -A mydb", "psql -F <str> -A <dbname>"},
		{"expanded quiet", "psql -x -q mydb", "psql -x -q <dbname>"},

		// Variable and pset flags
		{"variable", "psql -v name=value mydb", "psql -v <str> <dbname>"},
		{"pset", "psql -P format=unaligned mydb", "psql -P <str> <dbname>"},
		{"table attr", "psql -T 'border=1' mydb", "psql -T <str> <dbname>"},
		{"record separator", "psql -R '|' mydb", "psql -R <str> <dbname>"},

		// Long flags
		{"long host", "psql --host=db.example.com --port=5432 mydb", "psql --host=<host> --port=N <dbname>"},
		{"long dbname", "psql --dbname=mydb", "psql --dbname=<dbname>"},
		{"long username", "psql --username=admin", "psql --username=<user>"},
		{"long command", "psql --command='SELECT 1' mydb", "psql --command=<query> <dbname>"},
		{"long file", "psql --file=dump.sql mydb", "psql --file=<path> <dbname>"},
		{"long output", "psql --output=out.txt mydb", "psql --output=<path> <dbname>"},
		{"long log", "psql --log-file=log.txt mydb", "psql --log-file=<path> <dbname>"},

		// Multiple commands
		{"multiple commands", "psql -c 'CREATE TABLE t()' -c 'INSERT INTO t VALUES()' mydb", "psql -c <query> -c <query> <dbname>"},
		{"mixed command and file", "psql -c 'BEGIN' -f migrate.sql -c 'COMMIT' mydb", "psql -c <query> -f <path> -c <query> <dbname>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different databases collide", func(t *testing.T) {
		a := Normalize("psql -h localhost production")
		b := Normalize("psql -h localhost staging")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different queries collide", func(t *testing.T) {
		a := Normalize("psql -c 'SELECT * FROM users' mydb")
		b := Normalize("psql -c 'DELETE FROM orders' mydb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("psql mydb")
		subshell := Normalize("psql $(echo mydb)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
