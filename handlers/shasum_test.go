package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestShasum(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "shasum file.txt", "shasum <path>"},
		{"absolute path", "shasum /etc/passwd", "shasum <path>"},
		{"multiple files", "shasum file1.txt file2.txt file3.txt", "shasum <path>+"},

		// Algorithm flag (kept verbatim)
		{"algorithm 256", "shasum -a 256 file.txt", "shasum -a 256 <path>"},
		{"algorithm 512", "shasum -a 512 /tmp/data.bin", "shasum -a 512 <path>"},
		{"algorithm long", "shasum --algorithm 224 file.txt", "shasum --algorithm 224 <path>"},

		// Boolean flags
		{"check flag", "shasum -c checksums.txt", "shasum -c <path>"},
		{"binary flag", "shasum -b file.txt", "shasum -b <path>"},
		{"text flag", "shasum -t file.txt", "shasum -t <path>"},
		{"portable flag", "shasum -p file.txt", "shasum -p <path>"},

		// Combined flags
		{"algorithm and check", "shasum -a 256 -c sums.sha256", "shasum -a 256 -c <path>"},
		{"algorithm and binary", "shasum -a 512 -b largefile.bin", "shasum -a 512 -b <path>"},
		{"long flags", "shasum --algorithm 256 --binary file.txt", "shasum --algorithm 256 --binary <path>"},

		// Aliases
		{"sha256sum", "sha256sum file.txt", "sha256sum <path>"},
		{"sha1sum", "sha1sum file.txt", "sha1sum <path>"},
		{"sha512sum", "sha512sum data.bin", "sha512sum <path>"},
		{"md5sum", "md5sum file.txt", "md5sum <path>"},
		{"sha224sum", "sha224sum file.txt", "sha224sum <path>"},
		{"sha384sum", "sha384sum file.txt", "sha384sum <path>"},

		// Check mode with verify flags
		{"check with quiet", "shasum -c -q checksums.txt", "shasum -c -q <path>"},
		{"check with status", "shasum -c --status checksums.txt", "shasum -c --status <path>"},
		{"check with strict", "shasum -c --strict checksums.txt", "shasum -c --strict <path>"},
		{"check with warn", "shasum -c -w checksums.txt", "shasum -c -w <path>"},

		// Redirect
		{"with redirect", "shasum file.txt > sums.txt", "shasum <path> > <path>"},

		// No file (stdin)
		{"no file", "shasum -a 256", "shasum -a 256"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("shasum -a 256 /etc/passwd")
		b := shellshape.Normalize("shasum -a 256 /var/log/syslog")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different file counts collide", func(t *testing.T) {
		a := shellshape.Normalize("shasum file1.txt file2.txt")
		b := shellshape.Normalize("shasum a.txt b.txt c.txt d.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("shasum literal-file")
		subshell := shellshape.Normalize("shasum $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
