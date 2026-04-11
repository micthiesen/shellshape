package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUnzip(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"extract simple", "unzip archive.zip", "unzip <archive>"},
		{"list contents", "unzip -l archive.zip", "unzip -l <archive>"},
		{"test archive", "unzip -t data.zip", "unzip -t <archive>"},
		{"verbose list", "unzip -v backup.zip", "unzip -v <archive>"},
		{"pipe to stdout", "unzip -p fonts.zip", "unzip -p <archive>"},

		// Flags with arguments
		{"extract to dir", "unzip archive.zip -d /tmp/output", "unzip <archive> -d <path>"},
		{"password protected", "unzip -P secret archive.zip", "unzip -P <val> <archive>"},
		{"password and dir", "unzip -P mypass archive.zip -d /opt/dest", "unzip -P <val> <archive> -d <path>"},

		// Boolean modifier flags
		{"overwrite quietly", "unzip -oq archive.zip", "unzip -oq <archive>"},
		{"junk paths", "unzip -j archive.zip", "unzip -j <archive>"},
		{"never overwrite", "unzip -n archive.zip", "unzip -n <archive>"},
		{"case insensitive", "unzip -C archive.zip", "unzip -C <archive>"},

		// Extract specific files (member patterns)
		{"extract one file", "unzip archive.zip readme.txt", "unzip <archive> <pattern>"},
		{"extract multiple files", "unzip archive.zip file1.txt file2.txt file3.txt", "unzip <archive> <pattern>+"},
		{"extract glob", "unzip archive.zip '*.txt'", "unzip <archive> <pattern>"},

		// Exclude patterns with -x
		{"exclude files", "unzip archive.zip -x '*.log'", "unzip <archive> -x <pattern>"},
		{"exclude multiple", "unzip archive.zip -x '*.log' '*.tmp'", "unzip <archive> -x <pattern>+"},
		{"extract some exclude others", "unzip archive.zip '*.txt' -x readme.txt", "unzip <archive> <pattern> -x <pattern>"},

		// Combined usage
		{"full combo", "unzip -o archive.zip '*.conf' -x 'old/*' -d /etc/app", "unzip -o <archive> <pattern> -x <pattern> -d <path>"},
		{"flags before archive", "unzip -oq archive.zip -d /tmp", "unzip -oq <archive> -d <path>"},

		// Redirects
		{"with redirect", "unzip -l archive.zip > listing.txt", "unzip -l <archive> > <path>"},
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
		a := shellshape.Normalize("unzip project-v1.zip")
		b := shellshape.Normalize("unzip backup-2024.zip")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different destinations collide", func(t *testing.T) {
		a := shellshape.Normalize("unzip archive.zip -d /home/alice/docs")
		b := shellshape.Normalize("unzip archive.zip -d /var/data/output")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("unzip archive.zip")
		subshell := shellshape.Normalize("unzip $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
