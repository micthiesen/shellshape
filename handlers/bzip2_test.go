package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBzip2(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"compress file", "bzip2 file.txt", "bzip2 <path>"},
		{"decompress", "bzip2 -d file.bz2", "bzip2 -d <path>"},
		{"keep original", "bzip2 -k file.txt", "bzip2 -k <path>"},
		{"stdout", "bzip2 -c file.txt", "bzip2 -c <path>"},
		{"multiple files", "bzip2 a.txt b.txt c.txt", "bzip2 <path>+"},
		{"verbose", "bzip2 -v file.txt", "bzip2 -v <path>"},
		{"force", "bzip2 -f file.txt", "bzip2 -f <path>"},
		{"test", "bzip2 -t file.bz2", "bzip2 -t <path>"},

		// Compression levels (structural)
		{"level 1", "bzip2 -1 file.txt", "bzip2 -1 <path>"},
		{"level 9", "bzip2 -9 file.txt", "bzip2 -9 <path>"},

		// Combined flags
		{"decompress keep verbose", "bzip2 -dkv file.bz2", "bzip2 -dkv <path>"},
		{"force compress level", "bzip2 -f -9 file.txt", "bzip2 -f -9 <path>"},

		// Aliases
		{"bunzip2", "bunzip2 file.bz2", "bunzip2 <path>"},
		{"bzcat", "bzcat file.bz2", "bzcat <path>"},
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
		a := shellshape.Normalize("bzip2 important.txt")
		b := shellshape.Normalize("bzip2 secret.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("bzip2 file.txt")
		subshell := shellshape.Normalize("bzip2 $(echo file.txt)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
