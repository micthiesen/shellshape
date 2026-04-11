package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestW(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare w", "w", "w"},

		// Boolean flags
		{"no header", "w -h", "w -h"},
		{"sort by idle", "w -i", "w -i"},
		{"no dns", "w -n", "w -n"},
		{"combined flags", "w -hn", "w -hn"},

		// User positionals
		{"single user", "w root", "w <user>"},
		{"multiple users", "w root admin", "w <user>"},
		{"flags then user", "w -h root", "w -h <user>"},

		// Redirect
		{"with redirect", "w > /tmp/out", "w > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different users collide", func(t *testing.T) {
		a := shellshape.Normalize("w root")
		b := shellshape.Normalize("w admin")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("w root")
		subshell := shellshape.Normalize("w $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
