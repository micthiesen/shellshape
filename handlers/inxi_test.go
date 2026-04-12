package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestInxi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no args", "inxi", "inxi"},
		{"full info", "inxi -F", "inxi -F"},
		{"full anonymized", "inxi -Fz", "inxi -Fz"},
		{"basic", "inxi -b", "inxi -b"},
		{"graphics", "inxi -G", "inxi -G"},
		{"audio", "inxi -A", "inxi -A"},
		{"network", "inxi -N", "inxi -N"},
		{"disk", "inxi -D", "inxi -D"},
		{"system and machine", "inxi -S -M", "inxi -S -M"},
		{"color scheme", "inxi -c 5", "inxi -c N"},
		{"width", "inxi --width 120", "inxi --width N"},
		{"output file", "inxi -F --output-file /tmp/inxi.txt", "inxi -F --output-file <path>"},
		{"multiple sections", "inxi -G -A -N -D", "inxi -G -A -N -D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different color values collide", func(t *testing.T) {
		a := shellshape.Normalize("inxi -c 3")
		b := shellshape.Normalize("inxi -c 8")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("inxi -F")
		subshell := shellshape.Normalize("inxi $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
