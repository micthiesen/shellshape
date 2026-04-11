package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "du", "du"},
		{"single path", "du /home/user", "du <path>"},
		{"current dir", "du .", "du ."},
		{"human readable summary", "du -sh /var/log", "du -sh <path>"},
		{"all files human readable", "du -ah", "du -ah"},

		// Flags with arguments
		{"depth flag", "du -d 1 /var/log", "du -d N <path>"},
		{"max-depth long", "du --max-depth=2 /tmp", "du --max-depth=<val> <path>"},
		{"blocksize flag", "du -B 4096 /usr", "du -B N <path>"},
		{"ignore mask", "du -I '*.log' /var", "du -I <pattern> <path>"},
		{"threshold flag", "du -t 100M /usr", "du -t <threshold> <path>"},

		// Multiple positionals
		{"multiple paths", "du -sh /home /tmp /var", "du -sh <path>+"},
		{"summary with total", "du -sch /home/user /tmp", "du -sch <path>+"},

		// Edge cases
		{"bare name stays", "du mydir", "du mydir"},
		{"redirect", "du -sh /var > output.txt", "du -sh <path> > <path>"},
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
		a := shellshape.Normalize("du -sh /home/alice")
		b := shellshape.Normalize("du -sh /home/bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different depths collide", func(t *testing.T) {
		a := shellshape.Normalize("du -d 1 /var")
		b := shellshape.Normalize("du -d 3 /tmp")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("du -sh /home/user")
		subshell := shellshape.Normalize("du -sh $(pwd)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
