package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestAmixer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// List commands (no args)
		{"scontrols", "amixer scontrols", "amixer scontrols"},
		{"scontents", "amixer scontents", "amixer scontents"},
		{"controls", "amixer controls", "amixer controls"},
		{"contents", "amixer contents", "amixer contents"},

		// Simple set/get
		{"sset volume", "amixer sset Master 50%", "amixer sset <val>+"},
		{"sset with toggle", "amixer sset Master toggle", "amixer sset <val>+"},
		{"sset dB", "amixer sset Master 5dB+", "amixer sset <val>+"},
		{"sget", "amixer sget Master", "amixer sget <val>"},
		{"sset capture", "amixer sset Capture cap", "amixer sset <val>+"},

		// Control set/get
		{"cset", "amixer cset numid=3 1", "amixer cset <val>+"},
		{"cget", "amixer cget numid=3", "amixer cget <val>"},

		// Flags after subcommand
		{"card flag after", "amixer sset -c 0 Master 50%", "amixer sset -c N <val>+"},
		{"device flag after", "amixer sset -D pulse Master 80%", "amixer sset -D <val>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test
	t.Run("different controls collide", func(t *testing.T) {
		a := shellshape.Normalize("amixer sset Master 50%")
		b := shellshape.Normalize("amixer sset Capture 75%")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("amixer sset Master 50%")
		subshell := shellshape.Normalize("amixer sset $(evil-cmd) 50%")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
