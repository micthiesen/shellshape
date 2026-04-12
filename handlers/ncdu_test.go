package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNcdu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "ncdu", "ncdu"},
		{"single path", "ncdu /home", "ncdu <path>"},
		{"relative path", "ncdu ./projects", "ncdu <path>"},

		// Boolean flags
		{"one filesystem", "ncdu -x /", "ncdu -x <path>"},
		{"quiet", "ncdu -q /data", "ncdu -q <path>"},

		// Flags with arguments
		{"output file", "ncdu -o results.json", "ncdu -o <path>"},
		{"input file", "ncdu -f results.json", "ncdu -f <path>"},
		{"exclude pattern", "ncdu --exclude '*.txt' /home", "ncdu --exclude <pattern> <path>"},
		{"multiple excludes", "ncdu --exclude cache --exclude tmp /data", "ncdu --exclude <pattern> --exclude <pattern> <path>"},

		// Structural flags
		{"color dark", "ncdu --color dark /tmp", "ncdu --color dark <path>"},
		{"color off", "ncdu --color off /home", "ncdu --color off <path>"},

		// Combined
		{"complex", "ncdu -x -q --exclude '*.log' -o out.json /var", "ncdu -x -q --exclude <pattern> -o <path>+"},
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
		a := shellshape.Normalize("ncdu /home/alice")
		b := shellshape.Normalize("ncdu /var/log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ncdu /home")
		subshell := shellshape.Normalize("ncdu $(echo /home)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
