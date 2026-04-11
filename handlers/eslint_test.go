package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestEslint(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "eslint src/app.ts", "eslint <path>"},
		{"single directory", "eslint src/", "eslint <path>"},
		{"multiple paths", "eslint src/ lib/ tests/", "eslint <path>+"},

		// Boolean flags
		{"fix flag", "eslint --fix src/app.ts", "eslint --fix <path>"},
		{"cache flag", "eslint --cache src/", "eslint --cache <path>"},
		{"quiet flag", "eslint --quiet src/", "eslint --quiet <path>"},
		{"multiple boolean flags", "eslint --fix --cache --quiet src/", "eslint --fix --cache --quiet <path>"},

		// Flags with arguments
		{"config long", "eslint --config .eslintrc.json src/", "eslint --config <path>+"},
		{"config short", "eslint -c .eslintrc.json src/", "eslint -c <path>+"},
		{"ext flag", "eslint --ext .js,.ts src/", "eslint --ext <val> <path>"},
		{"rule flag", "eslint --rule 'no-console: error' src/", "eslint --rule <val> <path>"},
		{"ignore-path", "eslint --ignore-path .gitignore src/", "eslint --ignore-path <path>+"},
		{"output-file long", "eslint --output-file report.json src/", "eslint --output-file <path>+"},
		{"output-file short", "eslint -o report.json src/", "eslint -o <path>+"},
		{"format long", "eslint --format json src/", "eslint --format <val> <path>"},
		{"format short", "eslint -f stylish src/", "eslint -f <val> <path>"},
		{"max-warnings", "eslint --max-warnings 0 src/", "eslint --max-warnings N <path>"},
		{"cache-location", "eslint --cache --cache-location /tmp/.eslintcache src/", "eslint --cache --cache-location <path>+"},

		// Combined flags and multiple positionals
		{"complex invocation", "eslint --fix --ext .js,.ts --max-warnings 5 src/ lib/", "eslint --fix --ext <val> --max-warnings N <path>+"},
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
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("eslint --fix src/app.ts")
		b := shellshape.Normalize("eslint --fix lib/utils.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different rules collide", func(t *testing.T) {
		a := shellshape.Normalize("eslint --rule 'no-console: error' src/")
		b := shellshape.Normalize("eslint --rule 'no-unused-vars: warn' src/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("eslint src/app.ts")
		subshell := shellshape.Normalize("eslint $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
