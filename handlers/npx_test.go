package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNpx(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple command", "npx prettier --write src/", "npx prettier --write <path>"},
		{"scaffold", "npx create-react-app my-app", "npx create-react-app <val>"},
		{"command only", "npx cowsay", "npx cowsay"},
		{"multiple args", "npx cowsay hello world", "npx cowsay <val>+"},

		// Boolean flags
		{"yes flag", "npx -y ts-node script.ts", "npx -y ts-node <path>"},
		{"no flag", "npx --no cowsay hello", "npx --no cowsay <val>"},

		// Flags with arguments
		{"package flag short", "npx -p typescript tsc --noEmit", "npx -p <pkg> tsc --noEmit"},
		{"package flag long", "npx --package typescript tsc --noEmit", "npx --package <pkg> tsc --noEmit"},
		{"multiple packages", "npx -p @babel/core -p @babel/cli babel src/ -d lib/", "npx -p <pkg> -p <pkg> babel <path> -d <path>"},
		{"call flag short", "npx -c 'echo hello'", "npx -c <code>"},
		{"call flag long", "npx --call 'babel src/ -d lib/'", "npx --call <code>"},
		{"workspace flag", "npx -w packages/foo test", "npx -w <val> test"},

		// Bunx alias
		{"bunx alias", "bunx prettier --write .", "bunx prettier --write <val>"},

		// Scoped packages
		{"scoped package", "npx @angular/cli new my-project", "npx @angular/cli <val>+"},
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
	t.Run("different args collide", func(t *testing.T) {
		a := shellshape.Normalize("npx prettier --write src/app.ts")
		b := shellshape.Normalize("npx prettier --write lib/index.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different packages collide", func(t *testing.T) {
		a := shellshape.Normalize("npx -p typescript tsc")
		b := shellshape.Normalize("npx -p @types/node tsc")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("npx cowsay hello")
		subshell := shellshape.Normalize("npx cowsay $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
