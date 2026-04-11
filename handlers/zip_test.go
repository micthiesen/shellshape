package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestZip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple", "zip archive.zip file.txt", "zip <archive> <path>"},
		{"multiple files", "zip archive.zip file1.txt file2.txt file3.txt", "zip <archive> <path>+"},
		{"recursive", "zip -r archive.zip src/", "zip -r <archive> <path>"},
		{"recursive quiet", "zip -rq archive.zip dir/", "zip -rq <archive> <path>"},
		{"compress better", "zip -r -9 archive.zip src/", "zip -r -9 <archive> <path>"},
		{"delete entry", "zip -d archive.zip obsolete.txt", "zip -d <archive> <path>"},
		{"encrypt", "zip -e archive.zip secret.txt", "zip -e <archive> <path>"},
		{"update", "zip -u archive.zip changed.txt", "zip -u <archive> <path>"},

		// Flags with arguments
		{"password", "zip -P mysecret archive.zip file.txt", "zip -P <val> <archive> <path>"},
		{"temp path", "zip -b /tmp archive.zip files/", "zip -b <path> <archive> <path>"},
		{"date", "zip -t 01012024 archive.zip dir/", "zip -t <val> <archive> <path>"},
		{"suffixes", "zip -n .Z:.zip:.gz archive.zip dir/", "zip -n <val> <archive> <path>"},
		{"compression method", "zip -Z bzip2 archive.zip file.txt", "zip -Z <val> <archive> <path>"},

		// Exclude/include patterns
		{"exclude", "zip -r archive.zip dir/ -x '*.log'", "zip -r <archive> <path> -x <pattern>"},
		{"exclude multiple", "zip -r archive.zip dir/ -x '*.log' '*.tmp'", "zip -r <archive> <path> -x <pattern>+"},
		{"include", "zip -r archive.zip dir/ -i '*.go'", "zip -r <archive> <path> -i <pattern>"},

		// Multiple positionals with flags
		{"verbose multiple", "zip -rv archive.zip src/ lib/ bin/", "zip -rv <archive> <path>+"},

		// Stdin list
		{"stdin list", "zip -@ archive.zip", "zip -@ <archive>"},
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
	t.Run("different archives collide", func(t *testing.T) {
		a := shellshape.Normalize("zip -r project-v1.zip src/")
		b := shellshape.Normalize("zip -r backup-2024.zip src/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("zip archive.zip README.md config.yaml")
		b := shellshape.Normalize("zip archive.zip main.go utils.py")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("zip archive.zip file.txt")
		subshell := shellshape.Normalize("zip $(dangerous-command) file.txt")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
