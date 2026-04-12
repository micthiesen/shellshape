package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"find and replace", "sd foo bar", "sd <pattern> <val>"},
		{"with file", "sd foo bar file.txt", "sd <pattern> <val> <path>"},
		{"with multiple files", "sd foo bar a.txt b.txt c.txt", "sd <pattern> <val> <path>+"},
		{"regex pattern", `sd 'hello\s+world' 'hi'`, "sd <pattern> <val>"},

		// Boolean flags
		{"fixed strings", "sd -f hello world file.txt", "sd -f <pattern> <val> <path>"},
		{"fixed strings long", "sd --fixed-strings hello world", "sd --fixed-strings <pattern> <val>"},
		{"preview", "sd -p foo bar file.txt", "sd -p <pattern> <val> <path>"},
		{"string mode", "sd -s foo bar", "sd -s <pattern> <val>"},

		// Numeric flag
		{"max replacements", "sd -n 5 foo bar file.txt", "sd -n N <pattern> <val> <path>"},

		// Flags between positionals
		{"flag before pattern", "sd -f foo bar", "sd -f <pattern> <val>"},

		// Redirect
		{"with redirect", "sd foo bar < input.txt", "sd <pattern> <val> < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("sd foo bar file.txt")
		b := shellshape.Normalize("sd baz qux file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sd foo bar")
		subshell := shellshape.Normalize("sd $(dangerous-command) bar")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
