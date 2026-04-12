package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple pattern", "rg TODO", "rg <pattern>"},
		{"pattern and path", "rg todo src/", "rg <pattern> <path>"},
		{"pattern and multiple paths", "rg todo src/ lib/", "rg <pattern> <path>+"},

		// Boolean flags
		{"case insensitive", "rg -i pattern", "rg -i <pattern>"},
		{"word match", "rg -w TODO src/", "rg -w <pattern> <path>"},
		{"count", "rg -c pattern file.txt", "rg -c <pattern> <path>"},
		{"files with matches", "rg -l pattern", "rg -l <pattern>"},
		{"fixed strings", "rg -F 'literal.string'", "rg -F <pattern>"},

		// Flags with arguments
		{"type flag short", "rg -t py pattern", "rg -t <val> <pattern>"},
		{"type flag long", "rg --type rust pattern", "rg --type <val> <pattern>"},
		{"glob short", "rg -g '*.go' pattern", "rg -g <pattern>+"},
		{"glob long", "rg --glob '*.ts' pattern src/", "rg --glob <pattern>+ <path>"},
		{"explicit pattern", "rg -e foo -e bar", "rg -e <pattern> -e <pattern>"},
		{"max count", "rg --max-count 5 pattern", "rg --max-count N <pattern>"},

		// Context flags (numeric)
		{"after context", "rg -A 3 pattern file.txt", "rg -A N <pattern> <path>"},
		{"before context", "rg -B 5 pattern", "rg -B N <pattern>"},
		{"context", "rg -C 2 pattern", "rg -C N <pattern>"},
		{"fused context", "rg -A3 pattern file.txt", "rg -A N <pattern> <path>"},
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
		a := shellshape.Normalize("rg -i 'foo_bar' src/")
		b := shellshape.Normalize("rg -i 'baz_qux' lib/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different context numbers collide", func(t *testing.T) {
		a := shellshape.Normalize("rg -A 3 pattern file.go")
		b := shellshape.Normalize("rg -A 20 other file.py")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("rg pattern")
		subshell := shellshape.Normalize("rg $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
