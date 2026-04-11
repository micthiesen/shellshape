package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXxd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "xxd", "xxd"},
		{"input file", "xxd file.bin", "xxd <path>"},
		{"input and output file", "xxd input.bin output.hex", "xxd <path>+"},

		// Boolean flags
		{"reverse", "xxd -r dump.hex", "xxd -r <path>"},
		{"plain hex", "xxd -p file.bin", "xxd -p <path>"},
		{"c include", "xxd -i firmware.bin", "xxd -i <path>"},
		{"uppercase", "xxd -u file.bin", "xxd -u <path>"},
		{"binary bits", "xxd -b data.bin", "xxd -b <path>"},
		{"little endian", "xxd -e data.bin", "xxd -e <path>"},
		{"combined booleans", "xxd -r -p input.hex output.bin", "xxd -r -p <path>+"},

		// Flags with numeric arguments
		{"len flag", "xxd -l 120 file.bin", "xxd -l N <path>"},
		{"seek flag", "xxd -s 0x30 file.bin", "xxd -s N <path>"},
		{"cols flag", "xxd -c 20 file.bin", "xxd -c N <path>"},
		{"groupsize flag", "xxd -g 4 data.bin", "xxd -g N <path>"},
		{"offset flag", "xxd -o 0x100 file.bin", "xxd -o N <path>"},

		// Long flag forms
		{"long len", "xxd --len 256 file.bin", "xxd --len N <path>"},
		{"long cols", "xxd --cols 12 file.bin", "xxd --cols N <path>"},
		{"long groupsize", "xxd --groupsize 2 file.bin", "xxd --groupsize N <path>"},
		{"long seek", "xxd --seek 48 file.bin", "xxd --seek N <path>"},

		// Name flag (string value)
		{"name flag", "xxd -i -n my_array file.bin", "xxd -i -n <val> <path>"},
		{"long name flag", "xxd -i --name my_array file.bin", "xxd -i --name <val> <path>"},

		// Complex combinations
		{"len and cols and plain", "xxd -l 120 -ps -c 20 data.bin", "xxd -l N -ps -c N <path>"},
		{"seek negative", "xxd -s -0x30 file.bin", "xxd -s N <path>"},

		// Redirect
		{"with redirect", "xxd file.bin > output.hex", "xxd <path> > <path>"},

		// R flag (color mode)
		{"color flag", "xxd -R always file.bin", "xxd -R <val> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("xxd -l 120 secret.bin")
		b := shellshape.Normalize("xxd -l 256 firmware.dat")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric values collide", func(t *testing.T) {
		a := shellshape.Normalize("xxd -s 0x30 -c 16 file.bin")
		b := shellshape.Normalize("xxd -s 0xFF -c 32 data.bin")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xxd literal-arg")
		subshell := shellshape.Normalize("xxd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
