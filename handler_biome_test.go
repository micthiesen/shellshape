package shellshape

import "testing"

func TestBiome(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"check dir", "biome check src/", "biome check <path>"},
		{"lint file", "biome lint main.ts", "biome lint <path>"},
		{"format file", "biome format index.tsx", "biome format <path>"},
		{"ci subcommand", "biome ci .", "biome ci <path>"},
		{"migrate subcommand", "biome migrate", "biome migrate"},

		// Boolean flags
		{"write flag", "biome check --write src/", "biome check --write <path>"},
		{"unsafe flag", "biome lint --unsafe --write src/", "biome lint --unsafe --write <path>"},

		// Flags with arguments
		{"config-path", "biome format --config-path biome.json src/", "biome format --config-path <val> <path>"},
		{"colors flag", "biome check --colors force src/", "biome check --colors <val> <path>"},
		{"log-level flag", "biome lint --log-level info src/", "biome lint --log-level <val> <path>"},
		{"diagnostic-level flag", "biome ci --diagnostic-level error .", "biome ci --diagnostic-level <val> <path>"},

		// Multiple positionals
		{"multiple paths", "biome check src/ tests/ lib/", "biome check <path>+"},
		{"multiple files", "biome lint foo.ts bar.ts baz.ts", "biome lint <path>+"},

		// Mixed flags and paths
		{"mixed flags and paths", "biome check --write --colors off src/ tests/", "biome check --write --colors <val> <path>+"},
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
	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("biome check --write src/components/")
		b := Normalize("biome check --write src/utils/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different flag values collide", func(t *testing.T) {
		a := Normalize("biome lint --colors force src/")
		b := Normalize("biome lint --colors off src/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("biome check src/")
		subshell := Normalize("biome check $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
