package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestVitest(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "vitest", "vitest"},
		{"run subcommand", "vitest run", "vitest run"},
		{"bench subcommand", "vitest bench", "vitest bench"},
		{"watch subcommand", "vitest watch", "vitest watch"},

		// Positional test file patterns
		{"run with file", "vitest run src/utils.test.ts", "vitest run <path>"},
		{"bare with file", "vitest src/utils.test.ts", "vitest <path>"},
		{"multiple files", "vitest run src/a.test.ts src/b.test.ts", "vitest run <path>+"},

		// Flags with path arguments
		{"config long", "vitest --config ./vitest.config.ts", "vitest --config <path>"},
		{"config short", "vitest -c vitest.config.ts", "vitest -c <path>"},
		{"root flag", "vitest --root /home/user/project", "vitest --root <path>"},
		{"outputFile flag", "vitest --outputFile ./results.json", "vitest --outputFile <path>"},

		// Flags with value arguments
		{"reporter", "vitest --coverage --reporter json", "vitest --coverage --reporter <val>"},
		{"environment", "vitest run --environment jsdom src/", "vitest run --environment <val> <path>"},
		{"pool", "vitest bench --pool forks", "vitest bench --pool <val>"},

		// Flags with pattern arguments
		{"testNamePattern long", "vitest -t 'should handle errors'", "vitest -t <pattern>"},
		{"testNamePattern flag", "vitest --testNamePattern 'renders correctly'", "vitest --testNamePattern <pattern>"},
		{"grep flag", "vitest --grep 'auth.*login'", "vitest --grep <pattern>"},

		// Flags with numeric arguments
		{"bail", "vitest --bail 3", "vitest --bail N"},
		{"retry", "vitest --retry 2", "vitest --retry N"},
		{"maxConcurrency", "vitest --maxConcurrency 4", "vitest --maxConcurrency N"},

		// Boolean flags
		{"run flag", "vitest --run", "vitest --run"},
		{"coverage", "vitest --coverage", "vitest --coverage"},
		{"globals", "vitest --globals", "vitest --globals"},
		{"watch flag", "vitest --watch", "vitest --watch"},

		// Combined usage
		{"complex invocation", "vitest run --coverage --reporter verbose --bail 1 --environment jsdom src/components/", "vitest run --coverage --reporter <val> --bail N --environment <val> <path>"},
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
	t.Run("different test files collide", func(t *testing.T) {
		a := shellshape.Normalize("vitest run src/foo.test.ts")
		b := shellshape.Normalize("vitest run src/bar.test.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("vitest -t 'should render button'")
		b := shellshape.Normalize("vitest -t 'should handle click'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric values collide", func(t *testing.T) {
		a := shellshape.Normalize("vitest --bail 1 --retry 3")
		b := shellshape.Normalize("vitest --bail 5 --retry 10")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("vitest run src/utils.test.ts")
		subshell := shellshape.Normalize("vitest run $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
