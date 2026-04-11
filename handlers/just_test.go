package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestJust(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"recipe only", "just build", "just <recipe>"},
		{"recipe with args", "just deploy staging v1.2.3", "just <recipe> <arg>+"},
		{"no args", "just", "just"},
		{"list recipes", "just -l", "just -l"},
		{"dry run recipe", "just --dry-run test", "just --dry-run <recipe>"},
		{"init", "just --init", "just --init"},
		{"edit", "just -e", "just -e"},
		{"dump", "just --dump", "just --dump"},

		// Flags with path arguments
		{"justfile flag", "just -f Justfile.dev build", "just -f <path> <recipe>"},
		{"justfile long", "just --justfile /tmp/Justfile build", "just --justfile <path> <recipe>"},
		{"working directory", "just -d /home/user/project test", "just -d <path> <recipe>"},
		{"dotenv path", "just -E .env.local build", "just -E <path> <recipe>"},

		// Flags with value arguments
		{"shell flag", "just --shell bash build", "just --shell <val> <recipe>"},
		{"shell-arg flag", "just --shell-arg -cu build", "just --shell-arg <val> <recipe>"},
		{"color flag", "just --color always test", "just --color <val> <recipe>"},
		{"set variable", "just --set version 1.0 build", "just --set <val>+ <recipe>"},
		{"command flag", "just -c env", "just -c <val>"},
		{"completions", "just --completions bash", "just --completions <val>"},
		{"show recipe", "just -s deploy", "just -s <val>"},
		{"dump format", "just --dump-format json --dump", "just --dump-format <val> --dump"},
		{"group filter", "just --group deploy -l", "just --group <val> -l"},

		// Fused flags (--flag=value)
		{"justfile fused", "just --justfile=Justfile.dev build", "just --justfile=<path> <recipe>"},
		{"color fused", "just --color=auto test", "just --color=<val> <recipe>"},
		{"set fused", "just --set=version 1.0 build", "just --set=<val> <val> <recipe>"},

		// Multiple flags
		{"verbose dry-run", "just -v --dry-run test", "just -v --dry-run <recipe>"},
		{"justfile and working-dir", "just -f custom.just -d /tmp test arg1", "just -f <path> -d <path> <recipe> <arg>"},

		// Edge cases
		{"recipe with many args", "just deploy staging us-east-1 v2.0 --force", "just <recipe> <arg>+ --force"},
		{"quiet", "just -q build", "just -q <recipe>"},
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
	t.Run("different recipes collide", func(t *testing.T) {
		a := shellshape.Normalize("just build")
		b := shellshape.Normalize("just test")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different recipe args collide", func(t *testing.T) {
		a := shellshape.Normalize("just deploy staging")
		b := shellshape.Normalize("just deploy production")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different justfile paths collide", func(t *testing.T) {
		a := shellshape.Normalize("just -f Justfile.dev build")
		b := shellshape.Normalize("just -f Justfile.prod build")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("just build")
		subshell := shellshape.Normalize("just $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
