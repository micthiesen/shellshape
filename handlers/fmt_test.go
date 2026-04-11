package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFmt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare fmt", "fmt file.txt", "fmt <path>"},
		{"stdin no args", "fmt", "fmt"},
		{"format with width", "fmt -w 72 file.txt", "fmt -w N <path>"},

		// Boolean flags
		{"center", "fmt -c file.txt", "fmt -c <path>"},
		{"mail header", "fmt -m file.txt", "fmt -m <path>"},
		{"collapse spaces", "fmt -s file.txt", "fmt -s <path>"},
		{"multiple booleans", "fmt -c -s -p file.txt", "fmt -c -s -p <path>"},

		// Flags with arguments
		{"sentence chars", "fmt -d '.?!' file.txt", "fmt -d <str> <path>"},
		{"tab replace width", "fmt -l 4 file.txt", "fmt -l N <path>"},
		{"input tab stops", "fmt -t 4 file.txt", "fmt -t N <path>"},
		{"combined value flags", "fmt -l 4 -t 4 file.txt", "fmt -l N -t N <path>"},

		// Multiple positionals
		{"multiple files", "fmt file1.txt file2.txt", "fmt <path>+"},
		{"width with multiple files", "fmt -w 80 file1.txt file2.txt", "fmt -w N <path>+"},

		// Redirect
		{"with redirect", "fmt -w 72 file.txt > output.txt", "fmt -w N <path> > <path>"},
		{"stdin with redirect", "fmt -w 72 < input.txt", "fmt -w N < <path>"},
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
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("fmt readme.txt")
		b := shellshape.Normalize("fmt letter.md")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different widths collide", func(t *testing.T) {
		a := shellshape.Normalize("fmt -w 72 file.txt")
		b := shellshape.Normalize("fmt -w 80 file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fmt literal-arg")
		subshell := shellshape.Normalize("fmt $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
