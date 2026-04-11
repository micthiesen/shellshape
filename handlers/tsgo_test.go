package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTsgo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "tsgo", "tsgo"},
		{"version", "tsgo --version", "tsgo --version"},
		{"help", "tsgo --help", "tsgo --help"},
		{"noEmit only", "tsgo --noEmit", "tsgo --noEmit"},
		{"watch short", "tsgo -w", "tsgo -w"},
		{"build short", "tsgo -b", "tsgo -b"},

		// Flags with path arguments
		{"project short", "tsgo -p tsconfig.json", "tsgo -p <path>"},
		{"project long", "tsgo --project ./tsconfig.build.json", "tsgo --project <path>"},
		{"outDir", "tsgo --outDir dist", "tsgo --outDir <path>"},
		{"outFile", "tsgo --outFile bundle.js", "tsgo --outFile <path>"},
		{"rootDir", "tsgo --rootDir src", "tsgo --rootDir <path>"},
		{"declarationDir", "tsgo --declarationDir types", "tsgo --declarationDir <path>"},
		{"baseUrl", "tsgo --baseUrl ./src", "tsgo --baseUrl <path>"},

		// Value-consuming flags
		{"target", "tsgo --target ES2022 src/index.ts", "tsgo --target <val> <path>"},
		{"module", "tsgo --module commonjs src/index.ts", "tsgo --module <val> <path>"},
		{"jsx", "tsgo --jsx react-jsx src/App.tsx", "tsgo --jsx <val> <path>"},

		// Positional files
		{"single file", "tsgo src/index.ts", "tsgo <path>"},
		{"multiple files", "tsgo src/index.ts src/utils.ts src/types.ts", "tsgo <path>+"},

		// Combined flags and positionals
		{"project with noEmit", "tsgo -p tsconfig.json --noEmit", "tsgo -p <path> --noEmit"},
		{"project with watch", "tsgo -p tsconfig.json --noEmit --watch", "tsgo -p <path> --noEmit --watch"},
		{"outDir with strict", "tsgo --outDir dist --strict --noEmit", "tsgo --outDir <path> --strict --noEmit"},
		{"build with projects", "tsgo -b packages/core packages/utils", "tsgo -b <path>+"},
		{"strict with target", "tsgo --strict --target ES2022 --module ESNext src/index.ts", "tsgo --strict --target <val> --module <val> <path>"},

		// Long flag with = syntax
		{"target with equals", "tsgo --target=ES2022 src/index.ts", "tsgo --target=<val> <path>"},
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
	t.Run("different projects collide", func(t *testing.T) {
		a := shellshape.Normalize("tsgo -p tsconfig.json --noEmit")
		b := shellshape.Normalize("tsgo -p tsconfig.build.json --noEmit")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different target values collide", func(t *testing.T) {
		a := shellshape.Normalize("tsgo --target ES2020 src/index.ts")
		b := shellshape.Normalize("tsgo --target ES2022 src/main.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("tsgo src/index.ts")
		subshell := shellshape.Normalize("tsgo $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
