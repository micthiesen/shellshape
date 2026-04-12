package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXz(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"compress file", "xz file.txt", "xz <path>"},
		{"decompress", "xz -d file.xz", "xz -d <path>"},
		{"keep original", "xz -k file.txt", "xz -k <path>"},
		{"stdout", "xz -c file.txt", "xz -c <path>"},
		{"multiple files", "xz a.txt b.txt c.txt", "xz <path>+"},
		{"verbose", "xz -v file.txt", "xz -v <path>"},
		{"force", "xz -f file.txt", "xz -f <path>"},
		{"test", "xz -t file.xz", "xz -t <path>"},
		{"list", "xz -l file.xz", "xz -l <path>"},

		// Compression levels (structural, including -0)
		{"level 0", "xz -0 file.txt", "xz -0 <path>"},
		{"level 6", "xz -6 file.txt", "xz -6 <path>"},
		{"level 9", "xz -9 file.txt", "xz -9 <path>"},

		// Flags with arguments
		{"threads", "xz --threads 4 file.txt", "xz --threads N <path>"},
		{"memlimit", "xz --memlimit 256MiB file.txt", "xz --memlimit <val> <path>"},

		// Combined
		{"compress keep level", "xz -k -9 -v file.txt", "xz -k -9 -v <path>"},
		{"threaded compress", "xz --threads 8 -9 file.txt", "xz --threads N -9 <path>"},

		// Aliases
		{"unxz", "unxz file.xz", "unxz <path>"},
		{"xzcat", "xzcat file.xz", "xzcat <path>"},
		{"lzma", "lzma file.txt", "lzma <path>"},
		{"unlzma", "unlzma file.lzma", "unlzma <path>"},
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
		a := shellshape.Normalize("xz important.log")
		b := shellshape.Normalize("xz backup.tar")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xz file.txt")
		subshell := shellshape.Normalize("xz $(echo file.txt)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
