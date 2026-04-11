package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "ls", "ls"},
		{"single path", "ls /tmp", "ls <path>"},
		{"relative path", "ls ./src", "ls <path>"},
		{"file with extension", "ls file.txt", "ls <path>"},
		{"long listing", "ls -la", "ls -la"},
		{"long listing with path", "ls -la /home/user", "ls -la <path>"},
		{"recursive", "ls -lR src/ lib/", "ls -lR <path>+"},
		{"color flag with equals", "ls --color=always", "ls --color=<val>"},

		// Flags with arguments
		{"date format", "ls -D '%Y-%m-%d' /tmp", "ls -D <fmt> <path>"},

		// Multiple positionals
		{"multiple paths", "ls /tmp /var /etc", "ls <path>+"},
		{"flags between paths", "ls -l /tmp -a /var", "ls -l <path> -a <path>"},

		// Edge cases
		{"hidden files flag", "ls -A", "ls -A"},
		{"sort by size", "ls -lS", "ls -lS"},
		{"human readable", "ls -lh /usr/local/bin", "ls -lh <path>"},
		{"redirect", "ls -la /tmp > output.txt", "ls -la <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("ls -la /home/alice")
		b := shellshape.Normalize("ls -la /home/bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ls literal-arg")
		subshell := shellshape.Normalize("ls $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
