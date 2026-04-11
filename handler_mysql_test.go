package shellshape

import "testing"

func TestMysql(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"database only", "mysql mydb", "mysql <database>"},
		{"no args", "mysql", "mysql"},
		{"user and password prompt", "mysql -u root -p mydb", "mysql -u <user> -p <database>"},

		// Flags with arguments
		{"host flag", "mysql -h 10.0.0.1 mydb", "mysql -h <host> <database>"},
		{"long host flag", "mysql --host=db.example.com mydb", "mysql --host=<host> <database>"},
		{"port flag", "mysql -P 3307 mydb", "mysql -P N <database>"},
		{"execute query", "mysql -e 'SELECT * FROM users' mydb", "mysql -e <query> <database>"},
		{"long execute", "mysql --execute='SELECT 1' mydb", "mysql --execute=<query> <database>"},
		{"socket flag", "mysql -S /var/run/mysqld/mysqld.sock mydb", "mysql -S <path> <database>"},
		{"user long form", "mysql --user=admin mydb", "mysql --user=<user> <database>"},
		{"database flag", "mysql -D mydb", "mysql -D <database>"},
		{"defaults file", "mysql --defaults-file=/etc/my.cnf mydb", "mysql --defaults-file=<path> <database>"},

		// Full connection
		{"full connection", "mysql -h db.prod.internal -u deploy -p -P 3306 appdb", "mysql -h <host> -u <user> -p -P N <database>"},

		// Boolean flags
		{"batch mode", "mysql -B -e 'SELECT 1' mydb", "mysql -B -e <query> <database>"},
		{"skip column names", "mysql -N -B -e 'SHOW TABLES' mydb", "mysql -N -B -e <query> <database>"},
		{"verbose", "mysql -v mydb", "mysql -v <database>"},

		// Fused -p with password
		{"fused password", "mysql -u root -psecret mydb", "mysql -u <user> -p<str> <database>"},

		// Redirect input (restore from file)
		{"redirect input", "mysql -u root -p mydb < dump.sql", "mysql -u <user> -p <database> < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different databases collide", func(t *testing.T) {
		a := Normalize("mysql -u root -p production")
		b := Normalize("mysql -u root -p staging")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("mysql -h 10.0.0.1 -u root mydb")
		b := Normalize("mysql -h 10.0.0.2 -u admin mydb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("mysql mydb")
		subshell := Normalize("mysql $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
