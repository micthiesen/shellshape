package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDuf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "duf", "duf"},
		{"single path", "duf /home", "duf <path>"},
		{"multiple paths", "duf /home /var", "duf <path>+"},

		// Boolean flags
		{"all flag", "duf --all", "duf --all"},
		{"all short", "duf -a", "duf -a"},
		{"json", "duf --json", "duf --json"},

		// Value flags
		{"hide fstype", "duf --hide tmpfs", "duf --hide <val>"},
		{"only fstype", "duf --only local", "duf --only <val>"},
		{"output fields", "duf --output mountpoint,size,used", "duf --output <val>"},

		// Structural flags
		{"sort field", "duf --sort size", "duf --sort size"},
		{"theme", "duf --theme unicode", "duf --theme unicode"},

		// Combined
		{"complex", "duf --all --sort size --hide tmpfs /home", "duf --all --sort size --hide <val> <path>"},
		{"json output", "duf --json --only local /mnt", "duf --json --only <val> <path>"},
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
		a := shellshape.Normalize("duf /home/alice")
		b := shellshape.Normalize("duf /var/log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("duf /home")
		subshell := shellshape.Normalize("duf $(echo /home)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
