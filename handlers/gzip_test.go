package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGzip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"compress file", "gzip file.txt", "gzip <path>"},
		{"decompress file", "gzip -d archive.gz", "gzip -d <path>"},
		{"to stdout", "gzip -c file.txt", "gzip -c <path>"},
		{"keep original", "gzip -k file.txt", "gzip -k <path>"},

		// Compression levels
		{"fast compression", "gzip -1 file.txt", "gzip -1 <path>"},
		{"best compression", "gzip -9 file.txt", "gzip -9 <path>"},
		{"mid level", "gzip -5 largefile.log", "gzip -5 <path>"},
		{"long fast", "gzip --fast file.txt", "gzip --fast <path>"},
		{"long best", "gzip --best file.txt", "gzip --best <path>"},

		// Bundled flags
		{"bundled flags", "gzip -cd archive.gz", "gzip -cd <path>"},
		{"verbose keep", "gzip -vk file.txt", "gzip -vk <path>"},
		{"recursive verbose", "gzip -rv /var/log", "gzip -rv <path>"},

		// Suffix flag (consumes next token)
		{"custom suffix", "gzip -S .bak file.txt", "gzip -S <suffix> <path>"},
		{"long suffix", "gzip --suffix .backup file.txt", "gzip --suffix <suffix> <path>"},

		// Multiple positionals
		{"multiple files", "gzip file1.txt file2.txt file3.txt", "gzip <path>+"},
		{"decompress multiple", "gzip -d a.gz b.gz c.gz", "gzip -d <path>+"},

		// Testing and listing
		{"test integrity", "gzip -t archive.gz", "gzip -t <path>"},
		{"list info", "gzip -l archive.gz", "gzip -l <path>"},
		{"verbose list", "gzip -lv archive.gz", "gzip -lv <path>"},

		// Aliases
		{"gunzip", "gunzip archive.gz", "gunzip <path>"},
		{"zcat", "zcat archive.gz", "zcat <path>"},

		// Recursive
		{"recursive dir", "gzip -r /home/user/logs", "gzip -r <path>"},

		// With redirect
		{"redirect output", "gzip -c file.txt > output.gz", "gzip -c <path> > <path>"},
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
		a := shellshape.Normalize("gzip -9 /var/log/syslog")
		b := shellshape.Normalize("gzip -9 /tmp/data.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different suffixes collide", func(t *testing.T) {
		a := shellshape.Normalize("gzip -S .bak file.txt")
		b := shellshape.Normalize("gzip -S .old file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("gzip literal-arg")
		subshell := shellshape.Normalize("gzip $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
