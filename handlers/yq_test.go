package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestYq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"identity filter", "yq '.' file.yaml", "yq <filter> <path>"},
		{"filter only", "yq '.'", "yq <filter>"},
		{"no args", "yq", "yq"},
		{"raw output", "yq -r '.name' data.yaml", "yq -r <filter> <path>"},
		{"inplace edit", "yq -i '.name = \"test\"' file.yaml", "yq -i <filter> <path>"},
		{"null input", "yq -n '.a = 1'", "yq -n <filter>"},
		{"multiple input files", "yq '.key' a.yaml b.yaml c.yaml", "yq <filter> <path>+"},
		{"pretty print", "yq -P '.' data.yaml", "yq -P <filter> <path>"},
		{"colors", "yq -C '.' data.yaml", "yq -C <filter> <path>"},
		{"bundled flags", "yq -rC '.name' data.yaml", "yq -rC <filter> <path>"},

		// Flags with arguments
		{"output format short", "yq -o json '.' data.yaml", "yq -o <val> <filter> <path>"},
		{"output format long", "yq --output-format json '.' data.yaml", "yq --output-format <val> <filter> <path>"},
		{"input format", "yq -p json '.' data.json", "yq -p <val> <filter> <path>"},
		{"input format long", "yq --input-format xml '.' data.xml", "yq --input-format <val> <filter> <path>"},
		{"indent", "yq --indent 4 '.' data.yaml", "yq --indent N <filter> <path>"},
		{"indent short", "yq -I 4 '.' data.yaml", "yq -I N <filter> <path>"},

		// Subcommands (eval / eval-all)
		{"eval subcommand", "yq eval '.name' file.yaml", "yq eval <filter> <path>"},
		{"eval-all subcommand", "yq eval-all 'select(.name)' a.yaml b.yaml", "yq eval-all <filter> <path>+"},
		{"ea alias", "yq ea 'select(.name)' a.yaml b.yaml", "yq ea <filter> <path>+"},

		// Redirects
		{"with redirect", "yq '.' file.yaml > out.json", "yq <filter> <path> > <path>"},

		// Double dash
		{"double dash", "yq -- '.name' file.yaml", "yq -- <filter> <path>"},
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
	t.Run("different filters collide", func(t *testing.T) {
		a := shellshape.Normalize("yq '.name' data.yaml")
		b := shellshape.Normalize("yq '.age' data.yaml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("yq '.' config.yaml")
		b := shellshape.Normalize("yq '.' settings.yml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("yq '.name' file.yaml")
		subshell := shellshape.Normalize("yq $(dangerous-command) file.yaml")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
