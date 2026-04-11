package shellshape

import "testing"

func TestEnvsubst(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "envsubst", "envsubst"},
		{"with redirect in", "envsubst < template.txt", "envsubst < <path>"},
		{"with redirect in and out", "envsubst < template.txt > output.txt", "envsubst < <path> > <path>"},

		// Positional SHELL-FORMAT
		{"shell-format", "envsubst '$HOME $USER'", "envsubst <val>"},
		{"shell-format with redirect", "envsubst '$HOME' < t.txt > o.txt", "envsubst <val> < <path> > <path>"},

		// Flags
		{"variables flag short", "envsubst -v '$HOME $USER'", "envsubst -v <val>"},
		{"variables flag long", "envsubst --variables '$HOME $USER'", "envsubst --variables <val>"},
		{"version flag", "envsubst -V", "envsubst -V"},
		{"help flag", "envsubst --help", "envsubst --help"},

		// Pipeline usage
		{"in pipeline", "echo hello | envsubst", "echo <str> | envsubst"},
		{"pipeline with format", "cat t.txt | envsubst '$HOME'", "cat <path> | envsubst <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different shell-formats collide", func(t *testing.T) {
		a := Normalize("envsubst '$HOME $USER'")
		b := Normalize("envsubst '$PATH'")
		c := Normalize("envsubst '${DATABASE_URL}'")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("envsubst literal-arg")
		subshell := Normalize("envsubst $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
