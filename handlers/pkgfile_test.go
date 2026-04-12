package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPkgfile(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"search file", "pkgfile makepkg", "pkgfile <pattern>"},
		{"update", "pkgfile -u", "pkgfile -u"},
		{"list files", "pkgfile -l pacman", "pkgfile -l <val>"},
		{"list long", "pkgfile --list coreutils", "pkgfile --list <val>"},
		{"search binaries", "pkgfile -s -b gcc", "pkgfile -s -b <pattern>"},
		{"regex search", "pkgfile -r '.*bin.*'", "pkgfile -r <pattern>"},
		{"verbose", "pkgfile -v makepkg", "pkgfile -v <pattern>"},
		{"quiet glob", "pkgfile -q -g '*.conf'", "pkgfile -q -g <pattern>"},
		{"repo flag", "pkgfile -R extra makepkg", "pkgfile -R <val> <pattern>"},
		{"directories", "pkgfile -d /usr/bin", "pkgfile -d <pattern>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("pkgfile gcc")
		b := shellshape.Normalize("pkgfile vim")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pkgfile literal-arg")
		subshell := shellshape.Normalize("pkgfile $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
