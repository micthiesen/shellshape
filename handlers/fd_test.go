package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"pattern only", "fd TODO", "fd <pattern>"},
		{"pattern and path", "fd TODO src/", "fd <pattern> <path>"},
		{"no args", "fd", "fd"},

		// Boolean flags
		{"hidden", "fd -H pattern", "fd -H <pattern>"},
		{"no ignore", "fd -I pattern", "fd -I <pattern>"},
		{"unrestricted", "fd -u pattern", "fd -u <pattern>"},
		{"absolute paths", "fd -a pattern", "fd -a <pattern>"},
		{"list details", "fd -l pattern", "fd -l <pattern>"},

		// Type flag
		{"type file", "fd -t f pattern", "fd -t <val> <pattern>"},
		{"type dir long", "fd --type d pattern src/", "fd --type <val> <pattern> <path>"},

		// Extension flag
		{"extension short", "fd -e go pattern", "fd -e <val> <pattern>"},
		{"extension long", "fd --extension rs", "fd --extension <val>"},

		// Depth flags
		{"max depth short", "fd -d 3 pattern", "fd -d N <pattern>"},
		{"max depth long", "fd --max-depth 5 pattern", "fd --max-depth N <pattern>"},
		{"min depth", "fd --min-depth 2 pattern src/", "fd --min-depth N <pattern> <path>"},

		// Exclude flag
		{"exclude short", "fd -E '*.log' pattern", "fd -E <pattern>+"},
		{"exclude long", "fd --exclude node_modules pattern", "fd --exclude <pattern>+"},

		// Exec (rest of args)
		{"exec short", "fd -e go pattern -x wc -l", "fd -e <val> <pattern> -x <cmd...>"},
		{"exec long", "fd pattern --exec rm -f", "fd <pattern> --exec <cmd...>"},
		{"exec batch", "fd -e tmp -X rm", "fd -e <val> -X <cmd...>"},

		// Combinations
		{"hidden type depth", "fd -H -t f -d 2 pattern src/", "fd -H -t <val> -d N <pattern> <path>"},
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
		a := shellshape.Normalize("fd 'foo_bar' src/")
		b := shellshape.Normalize("fd 'baz_qux' lib/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fd pattern")
		subshell := shellshape.Normalize("fd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
