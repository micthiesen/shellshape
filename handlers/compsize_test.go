package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestCompsize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single path", "compsize /mnt/data", "compsize <path>"},
		{"multiple paths", "compsize /mnt/data /mnt/games", "compsize <path>+"},
		{"relative path", "compsize ./subvol", "compsize <path>"},
		{"home path", "compsize ~/Documents", "compsize <path>"},
		{"current dir", "compsize .", "compsize <path>"},

		// Flags
		{"one filesystem", "compsize -x /mnt/data", "compsize -x <path>"},
		{"flag after path", "compsize /mnt/data -x", "compsize <path> -x"},

		// Multiple paths with flag
		{"flag multi paths", "compsize -x /mnt/data /mnt/games /home", "compsize -x <path>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("compsize /mnt/data")
		b := shellshape.Normalize("compsize /mnt/games")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("compsize /mnt/data")
		subshell := shellshape.Normalize("compsize $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
