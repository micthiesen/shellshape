package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestShellshape(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare shellshape", "shellshape", "shellshape"},
		{"single word command", "shellshape ls", "shellshape <command>"},
		{"simple command", "shellshape echo hello", "shellshape <command>"},
		{"command with flags", "shellshape grep -r foo /bar", "shellshape <command>"},
		{"quoted command", `shellshape "git commit -m fix"`, "shellshape <command>"},
		{"complex command", "shellshape curl -s -H 'Content-Type: application/json' https://example.com", "shellshape <command>"},
		{"command with pipes as single arg", `shellshape 'echo hello | grep h'`, "shellshape <command>"},
		{"multiple words", "shellshape docker run -it ubuntu bash", "shellshape <command>"},

		// List subcommand
		{"list subcommand", "shellshape list", "shellshape list"},

		// With redirect
		{"with redirect", "shellshape echo hi > /tmp/out.txt", "shellshape <command> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values -> same shape
	t.Run("different commands collide", func(t *testing.T) {
		a := shellshape.Normalize("shellshape echo hello")
		b := shellshape.Normalize("shellshape grep -r foo /bar")
		c := shellshape.Normalize("shellshape curl -s https://example.com")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("shellshape literal-arg")
		subshell := shellshape.Normalize("shellshape $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
