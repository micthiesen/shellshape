package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestAwk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"print field", "awk '{print $1}' file.txt", "awk <awk-prog> <path>"},
		{"pattern match", "awk '/error/' log.txt", "awk <awk-prog> <path>"},
		{"program only", "awk 'BEGIN{print \"hello\"}'", "awk <awk-prog>"},
		{"multiple files", "awk 'NR==1' file1.txt file2.txt file3.txt", "awk <awk-prog> <path>+"},

		// Flags with arguments
		{"field sep", "awk -F ',' '{print $1}' data.csv", "awk -F <val> <awk-prog> <path>"},
		{"field sep colon", "awk -F: '{print $1}' /etc/passwd", "awk -F <val> <awk-prog> <path>"},
		{"variable assignment", "awk -v OFS='\\t' '{print $1,$2}' data.txt", "awk -v <val> <awk-prog> <path>"},
		{"multiple vars", "awk -v x=1 -v y=2 '{print x,y}' file.txt", "awk -v <val> -v <val> <awk-prog> <path>"},
		{"program file", "awk -f script.awk data.txt", "awk -f <path>+"},
		{"program file no input", "awk -f prog.awk", "awk -f <path>"},

		// Redirect
		{"with redirect", "awk '{print $1}' input.txt > output.txt", "awk <awk-prog> <path> > <path>"},

		// Edge: -f means no positional program
		{"f flag with multiple inputs", "awk -f prog.awk a.txt b.txt", "awk -f <path>+"},
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
	t.Run("different programs collide", func(t *testing.T) {
		a := shellshape.Normalize("awk '{print $1}' data.csv")
		b := shellshape.Normalize("awk '{print $NF}' data.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different field separators collide", func(t *testing.T) {
		a := shellshape.Normalize("awk -F ',' '{print $1}' file.csv")
		b := shellshape.Normalize("awk -F ':' '{print $1}' file.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("awk literal-arg")
		subshell := shellshape.Normalize("awk $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
