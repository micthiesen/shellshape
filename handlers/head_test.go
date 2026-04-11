package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestHead(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "head file.txt", "head <path>"},
		{"no args", "head", "head"},
		{"stdin implicit", "head -n 5", "head -n N"},

		// Flags with numeric arguments
		{"-n lines", "head -n 10 file.txt", "head -n N <path>"},
		{"-c bytes", "head -c 100 file.txt", "head -c N <path>"},
		{"--lines long", "head --lines 20 file.txt", "head --lines N <path>"},
		{"--bytes long", "head --bytes 512 file.txt", "head --bytes N <path>"},
		{"--lines=val", "head --lines=50 file.txt", "head --lines=<val> <path>"},

		// Legacy -N form (handled by classifyToken as number)
		{"legacy count", "head -20 file.txt", "head N <path>"},

		// Multiple files
		{"multiple files", "head -n 5 a.txt b.txt c.txt", "head -n N <path>+"},

		// With redirect
		{"redirect", "head -n 10 file.txt > out.txt", "head -n N <path> > <path>"},
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
		a := shellshape.Normalize("head -n 10 server.log")
		b := shellshape.Normalize("head -n 10 access.log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different counts collide", func(t *testing.T) {
		a := shellshape.Normalize("head -n 5 file.txt")
		b := shellshape.Normalize("head -n 100 file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("head literal-arg")
		subshell := shellshape.Normalize("head $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
