package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWhich(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"which simple", "which go", "which <cmd>"},
		{"which no args", "which", "which"},
		{"which multiple", "which go python node", "which <cmd>"},
		{"which -a flag", "which -a go", "which -a <cmd>"},
		{"which -s flag", "which -s go", "which -s <cmd>"},

		// whereis
		{"whereis simple", "whereis python", "whereis <cmd>"},
		{"whereis -b flag", "whereis -b python", "whereis -b <cmd>"},
		{"whereis -B path", "whereis -B /usr/bin python", "whereis -B <path> <cmd>"},
		{"whereis -M path", "whereis -M /usr/share/man python", "whereis -M <path> <cmd>"},
		{"whereis -S path", "whereis -S /usr/src python", "whereis -S <path> <cmd>"},
		{"whereis mixed flags", "whereis -b -m python", "whereis -b -m <cmd>"},

		// type
		{"type simple", "type ls", "type <cmd>"},
		{"type -t flag", "type -t ls", "type -t <cmd>"},
		{"type -a flag", "type -a ls", "type -a <cmd>"},
		{"type -p flag", "type -p ls", "type -p <cmd>"},
		{"type multiple", "type ls cat grep", "type <cmd>"},

		// command
		{"command -v", "command -v node", "command -v <cmd>"},
		{"command -V", "command -V node", "command -V <cmd>"},
		{"command -p", "command -p ls", "command -p <cmd>"},

		// Redirect
		{"which with redirect", "which go > /tmp/out", "which <cmd> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different command names → same shape
	t.Run("different commands collide", func(t *testing.T) {
		a := shellshape.Normalize("which kdotool")
		b := shellshape.Normalize("which xdg-close")
		c := shellshape.Normalize("which ydotool")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different type lookups collide", func(t *testing.T) {
		a := shellshape.Normalize("type ls")
		b := shellshape.Normalize("type grep")
		if a != b {
			t.Errorf("expected same shape: %q, %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("which literal-arg")
		subshell := shellshape.Normalize("which $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
