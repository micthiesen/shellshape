package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare pipe", "pv", "pv"},
		{"single file", "pv file.iso", "pv <path>"},
		{"multiple files", "pv a.bin b.bin c.bin", "pv <path>+"},
		{"absolute path", "pv /dev/sda", "pv <path>"},

		// Boolean flags
		{"progress flags", "pv -p -t -e -r file.bin", "pv -p -t -e -r <path>"},
		{"line mode", "pv -l file.txt", "pv -l <path>"},
		{"force flag", "pv -f file.iso", "pv -f <path>"},

		// Numeric flags
		{"size flag", "pv -s 1024m file.bin", "pv -s N <path>"},
		{"size flag long", "pv --size 500m file.bin", "pv --size N <path>"},
		{"rate limit", "pv -L 10m file.dat", "pv -L N <path>"},
		{"rate limit long", "pv --rate-limit 1m file.dat", "pv --rate-limit N <path>"},

		// Value flags
		{"name flag", "pv -N copying file.iso", "pv -N <val> <path>"},
		{"name flag long", "pv --name transfer file.iso", "pv --name <val> <path>"},

		// Combined
		{"typical usage", "pv -s 4g -N backup disk.img", "pv -s N -N <val> <path>"},

		// Redirect
		{"with redirect", "pv file.iso > /dev/sdb", "pv <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("pv -s 100m alpha.bin")
		b := shellshape.Normalize("pv -s 500m beta.bin")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pv file.bin")
		subshell := shellshape.Normalize("pv $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
