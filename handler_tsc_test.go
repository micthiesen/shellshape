package shellshape

import "testing"

func TestTsc(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "tsc", "tsc"},
		{"single file", "tsc src/index.ts", "tsc <path>"},
		{"multiple files", "tsc src/a.ts src/b.ts src/c.ts", "tsc <path>+"},

		// Project flag
		{"project short", "tsc -p tsconfig.json", "tsc -p <path>"},
		{"project long", "tsc --project tsconfig.build.json", "tsc --project <path>"},

		// Path-consuming flags
		{"outDir", "tsc --outDir dist", "tsc --outDir <path>"},
		{"rootDir", "tsc --rootDir src", "tsc --rootDir <path>"},
		{"outDir and rootDir", "tsc --outDir dist --rootDir src", "tsc --outDir <path> --rootDir <path>"},
		{"outFile", "tsc --outFile bundle.js", "tsc --outFile <path>"},
		{"declarationDir", "tsc --declarationDir types", "tsc --declarationDir <path>"},
		{"baseUrl", "tsc --baseUrl ./src", "tsc --baseUrl <path>"},

		// Value-consuming flags
		{"target", "tsc --target ES2022 src/index.ts", "tsc --target <val> <path>"},
		{"module", "tsc --module commonjs src/index.ts", "tsc --module <val> <path>"},
		{"moduleResolution", "tsc --moduleResolution node src/index.ts", "tsc --moduleResolution <val> <path>"},
		{"jsx", "tsc --jsx react-jsx src/App.tsx", "tsc --jsx <val> <path>"},
		{"lib", "tsc --lib ES2022 src/index.ts", "tsc --lib <val> <path>"},

		// Boolean flags
		{"noEmit", "tsc --noEmit", "tsc --noEmit"},
		{"strict", "tsc --strict src/index.ts", "tsc --strict <path>"},
		{"declaration", "tsc -d", "tsc -d"},
		{"watch short", "tsc -w", "tsc -w"},
		{"build short", "tsc -b", "tsc -b"},
		{"incremental", "tsc --incremental", "tsc --incremental"},

		// Combined flags
		{"declaration and outDir", "tsc --declaration --emitDeclarationOnly --outDir types", "tsc --declaration --emitDeclarationOnly --outDir <path>"},
		{"strict with target", "tsc --strict --target ES2022 --module ESNext src/index.ts", "tsc --strict --target <val> --module <val> <path>"},
		{"build with projects", "tsc -b packages/core packages/utils", "tsc -b <path>+"},
		{"watch with project", "tsc -w -p tsconfig.json", "tsc -w -p <path>"},

		// Long flag with =
		{"target with equals", "tsc --target=ES2022 src/index.ts", "tsc --target=<val> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("tsc src/foo.ts")
		b := Normalize("tsc src/bar.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different target values collide", func(t *testing.T) {
		a := Normalize("tsc --target ES2020 src/index.ts")
		b := Normalize("tsc --target ES2022 src/main.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tsc src/index.ts")
		subshell := Normalize("tsc $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
