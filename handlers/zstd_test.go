package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestZstd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"compress file", "zstd file.txt", "zstd <path>"},
		{"decompress file", "zstd -d archive.zst", "zstd -d <path>"},
		{"to stdout", "zstd -c file.txt", "zstd -c <path>"},
		{"keep original", "zstd -k file.txt", "zstd -k <path>"},
		{"force overwrite", "zstd -f file.txt", "zstd -f <path>"},
		{"verbose", "zstd -v file.txt", "zstd -v <path>"},
		{"test integrity", "zstd -t archive.zst", "zstd -t <path>"},
		{"quiet mode", "zstd -q file.txt", "zstd -q <path>"},
		{"remove source", "zstd --rm file.txt", "zstd --rm <path>"},
		{"ultra mode", "zstd --ultra -19 file.txt", "zstd --ultra -19 <path>"},

		// Compression levels
		{"level 1", "zstd -1 file.txt", "zstd -1 <path>"},
		{"level 9", "zstd -9 file.txt", "zstd -9 <path>"},
		{"level 19", "zstd -19 file.txt", "zstd -19 <path>"},

		// Flags with arguments
		{"output file", "zstd -o output.zst input.txt", "zstd -o <path>+"},
		{"dictionary", "zstd -D dict.bin input.txt", "zstd -D <path>+"},
		{"threads", "zstd --threads 4 file.txt", "zstd --threads N <path>"},
		{"fast mode fused", "zstd --fast=3 file.txt", "zstd --fast=N <path>"},
		{"train flag", "zstd --train samples/", "zstd --train <path>"},

		// Bundled flags
		{"bundled decompress keep", "zstd -dk archive.zst", "zstd -dk <path>"},
		{"bundled force verbose", "zstd -fv file.txt", "zstd -fv <path>"},

		// Multiple positionals
		{"multiple files", "zstd file1.txt file2.txt file3.txt", "zstd <path>+"},
		{"decompress multiple", "zstd -d a.zst b.zst c.zst", "zstd -d <path>+"},

		// Aliases
		{"unzstd", "unzstd archive.zst", "unzstd <path>"},
		{"zstdcat", "zstdcat archive.zst", "zstdcat <path>"},
		{"zstdmt", "zstdmt -19 file.txt", "zstdmt -19 <path>"},

		// With redirect
		{"redirect output", "zstd -c file.txt > output.zst", "zstd -c <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("zstd -19 /var/log/syslog")
		b := shellshape.Normalize("zstd -19 /tmp/data.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output paths collide", func(t *testing.T) {
		a := shellshape.Normalize("zstd -o /tmp/out1.zst input.txt")
		b := shellshape.Normalize("zstd -o /tmp/out2.zst input.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("zstd literal-arg")
		subshell := shellshape.Normalize("zstd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
