package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSqlite3(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"db and query", `sqlite3 mydata.db "SELECT * FROM memos"`, "sqlite3 <path> <sql>"},
		{"db only", "sqlite3 mydata.db", "sqlite3 <path>"},
		{"db with flags", "sqlite3 -header -column mydata.db 'SELECT * FROM users'", "sqlite3 -header -column <path> <sql>"},
		{"db path absolute", `sqlite3 /tmp/test.db "SELECT 1"`, "sqlite3 <path> <sql>"},
		{"db path tilde", `sqlite3 ~/.local/share/app/data.db "SELECT count(*) FROM events"`, "sqlite3 <path> <sql>"},
		{"line mode", `sqlite3 -line mydata.db "SELECT * FROM memos WHERE priority > 20"`, "sqlite3 -line <path> <sql>"},
		{"separator flag", `sqlite3 -separator "," data.db "SELECT * FROM t"`, "sqlite3 -separator <val> <path> <sql>"},
		{"bare command", "sqlite3", "sqlite3"},
		{"with redirect", `sqlite3 data.db "SELECT 1" > out.txt`, "sqlite3 <path> <sql> > <path>"},
		{"in pipeline", `sqlite3 data.db "SELECT name FROM t" | sort`, "sqlite3 <path> <sql> | sort"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different queries collide", func(t *testing.T) {
		a := shellshape.Normalize(`sqlite3 data.db "SELECT * FROM users"`)
		b := shellshape.Normalize(`sqlite3 data.db "DELETE FROM sessions WHERE expired = 1"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different databases collide", func(t *testing.T) {
		a := shellshape.Normalize(`sqlite3 /tmp/a.db "SELECT 1"`)
		b := shellshape.Normalize(`sqlite3 /tmp/b.db "SELECT 1"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize(`sqlite3 data.db "SELECT 1"`)
		subshell := shellshape.Normalize(`sqlite3 data.db $(cat query.sql)`)
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
