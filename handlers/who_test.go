package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWho(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare who", "who", "who"},
		{"who am i", "who am i", "who am i"},

		// Boolean flags
		{"all info", "who -a", "who -a"},
		{"with header", "who -H -u", "who -H -u"},
		{"quick mode", "who -q", "who -q"},
		{"multiple flags", "who -bTu", "who -bTu"},

		// File positional
		{"alternate file", "who /var/log/utx.log", "who <path>"},
		{"flags then file", "who -H /var/log/utx.log", "who -H <path>"},

		// Redirect
		{"with redirect", "who > /tmp/out", "who > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("who /var/log/utx.log")
		b := shellshape.Normalize("who /var/run/utmpx")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("who /var/log/utx.log")
		subshell := shellshape.Normalize("who $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
