package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMan(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple page", "man ls", "man ls"},
		{"section and page", "man 3 printf", "man 3 printf"},
		{"section 2", "man 2 stat", "man 2 stat"},

		// Boolean flags
		{"all pages", "man -a stat", "man -a stat"},
		{"whatis mode", "man -f printf", "man -f printf"},
		{"apropos mode", "man -k copy", "man -k copy"},
		{"show path", "man -w ls", "man -w ls"},
		{"full text search", "man -K printf", "man -K printf"},

		// Flags with arguments
		{"manpath flag", "man -M /usr/local/man ls", "man -M <path> ls"},
		{"pager flag", "man -P less ls", "man -P <val> ls"},
		{"sections flag", "man -S 1:8 ls", "man -S <val> ls"},
		{"section via -s", "man -s 3 printf", "man -s <val> printf"},
		{"arch flag", "man -m amd64 ls", "man -m <val> ls"},
		{"preprocessor flag", "man -p tev ls", "man -p <val> ls"},

		// Combined flags
		{"apropos with section", "man -k -S 1:3 socket", "man -k -S <val> socket"},
		{"all with section prefix", "man -a 3 printf", "man -a 3 printf"},

		// Multiple pages
		{"multiple pages", "man ls cat grep", "man ls cat grep"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different manpaths collide", func(t *testing.T) {
		a := shellshape.Normalize("man -M /usr/local/man ls")
		b := shellshape.Normalize("man -M /opt/man ls")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pagers collide", func(t *testing.T) {
		a := shellshape.Normalize("man -P less ls")
		b := shellshape.Normalize("man -P more ls")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("man ls")
		subshell := shellshape.Normalize("man $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
