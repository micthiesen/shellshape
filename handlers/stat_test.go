package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestStat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple file", "stat file.txt", "stat <path>"},
		{"absolute path", "stat /etc/hosts", "stat <path>"},
		{"multiple files", "stat file1.txt file2.txt file3.txt", "stat <path>+"},

		// Boolean flags
		{"dereference", "stat -L /etc/hosts", "stat -L <path>"},
		{"verbose format", "stat -x /tmp/foo", "stat -x <path>"},
		{"combined booleans", "stat -Lx /tmp/foo", "stat -Lx <path>"},
		{"shell output", "stat -s myfile.go", "stat -s <path>"},

		// Flags with format arguments
		{"format flag -f", "stat -f '%z' myfile.txt", "stat -f <fmt> <path>"},
		{"time format -t", "stat -t '%Y%m%d' file.log", "stat -t <fmt> <path>"},
		{"linux format -c", "stat -c '%a' /etc/passwd", "stat -c <fmt> <path>"},
		{"combined format and time", "stat -f '%Sm' -t '%Y%m%d' file.log", "stat -f <fmt> -t <fmt> <path>"},
		{"long format with equals", "stat --format='%s' data.tar", "stat --format=<val> <path>"},
		{"long printf with equals", "stat --printf='%s' data.tar", "stat --printf=<val> <path>"},
		{"long format space", "stat --format '%n' data.tar", "stat --format <fmt> <path>"},
		{"long printf space", "stat --printf '%n' data.tar", "stat --printf <fmt> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("stat /etc/passwd")
		b := shellshape.Normalize("stat /var/log/syslog")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different format strings collide", func(t *testing.T) {
		a := shellshape.Normalize("stat -f '%z' file1.txt")
		b := shellshape.Normalize("stat -f '%Sm' file2.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("stat literal-arg")
		subshell := shellshape.Normalize("stat $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
