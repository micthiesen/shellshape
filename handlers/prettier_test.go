package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPrettier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"write single file", "prettier --write src/app.tsx", "prettier --write <path>"},
		{"check directory", "prettier --check src/", "prettier --check <path>"},
		{"short write flag", "prettier -w file.js", "prettier -w <path>"},
		{"short check flag", "prettier -c file.ts", "prettier -c <path>"},

		// Flags with arguments
		{"config path", "prettier --config .prettierrc --write .", "prettier --config <path> --write ."},
		{"ignore-path", "prettier --ignore-path .gitignore --write src/", "prettier --ignore-path <path> --write <path>"},
		{"parser value", "prettier --parser typescript --write src/", "prettier --parser <val> --write <path>"},
		{"tab-width numeric", "prettier --tab-width 4 file.js", "prettier --tab-width N <path>"},
		{"print-width numeric", "prettier --print-width 100 file.ts", "prettier --print-width N <path>"},
		{"trailing-comma value", "prettier --trailing-comma all --write .", "prettier --trailing-comma <val> --write ."},
		{"arrow-parens value", "prettier --arrow-parens avoid --write .", "prettier --arrow-parens <val> --write ."},
		{"end-of-line value", "prettier --end-of-line lf --write src/", "prettier --end-of-line <val> --write <path>"},
		{"log-level value", "prettier --log-level warn --write .", "prettier --log-level <val> --write ."},

		// Multiple positionals
		{"multiple files", "prettier --write src/app.tsx src/index.ts lib/utils.js", "prettier --write <path>+"},

		// Boolean flags
		{"single-quote", "prettier --single-quote --write file.js", "prettier --single-quote --write <path>"},
		{"no-semi", "prettier --no-semi --write file.ts", "prettier --no-semi --write <path>"},
		{"list-different", "prettier -l src/", "prettier -l <path>"},
		{"no-config", "prettier --no-config --write file.js", "prettier --no-config --write <path>"},

		// Combined options
		{"full config", "prettier --single-quote --trailing-comma all --tab-width 4 --print-width 100 --write src/", "prettier --single-quote --trailing-comma <val> --tab-width N --print-width N --write <path>"},

		// Edge cases
		{"plugin flag", "prettier --plugin prettier-plugin-tailwindcss --write .", "prettier --plugin <val> --write ."},
		{"range flags", "prettier --range-start 0 --range-end 100 file.js", "prettier --range-start N --range-end N <path>"},
		{"cache-location", "prettier --cache --cache-location /tmp/prettier-cache --write .", "prettier --cache --cache-location <path> --write ."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("prettier --write src/app.tsx")
		b := shellshape.Normalize("prettier --write lib/utils.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different parser values collide", func(t *testing.T) {
		a := shellshape.Normalize("prettier --parser typescript --write src/")
		b := shellshape.Normalize("prettier --parser babel --write lib/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("prettier --write file.js")
		subshell := shellshape.Normalize("prettier --write $(find . -name '*.js')")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
