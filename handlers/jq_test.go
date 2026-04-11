package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestJq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"identity filter", "jq '.' file.json", "jq <filter> <path>"},
		{"filter only", "jq '.'", "jq <filter>"},
		{"raw output", "jq -r '.name' data.json", "jq -r <filter> <path>"},
		{"compact output", "jq -c '.[] | .id' data.json", "jq -c <filter> <path>"},
		{"null input", "jq -n '{a: 1}'", "jq -n <filter>"},
		{"multiple input files", "jq '.key' a.json b.json c.json", "jq <filter> <path>+"},
		{"bundled flags", "jq -rS '.name' data.json", "jq -rS <filter> <path>"},

		// Flags with arguments
		{"from-file", "jq -f script.jq data.json", "jq -f <path>+"},
		{"from-file long", "jq --from-file script.jq data.json", "jq --from-file <path>+"},
		{"indent", "jq --indent 4 '.' data.json", "jq --indent N <filter> <path>"},
		{"library path", "jq -L /usr/lib/jq '.' data.json", "jq -L <path> <filter> <path>"},
		{"arg flag", "jq --arg name value '.[$name]' data.json", "jq --arg <name> <val> <filter> <path>"},
		{"argjson flag", "jq --argjson count 42 '. + {count: $count}' data.json", "jq --argjson <name> <val> <filter> <path>"},
		{"slurpfile flag", "jq --slurpfile vars vars.json '.+$vars' data.json", "jq --slurpfile <name> <path> <filter> <path>"},
		{"rawfile flag", "jq --rawfile tmpl template.txt '.+$tmpl' data.json", "jq --rawfile <name> <path> <filter> <path>"},

		// --args and --jsonargs make remaining positionals into <val>
		{"args flag", "jq -n --args '$ARGS.positional' foo bar baz", "jq -n --args <filter> <val>+"},
		{"jsonargs flag", "jq -n --jsonargs '$ARGS.positional' 1 2 3", "jq -n --jsonargs <filter> <val>+"},

		// Redirects
		{"with redirect", "jq '.' file.json > out.json", "jq <filter> <path> > <path>"},

		// Edge: -- separator
		{"double dash", "jq -- '.name' file.json", "jq -- <filter> <path>"},
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
		a := shellshape.Normalize("jq '.name' data.json")
		b := shellshape.Normalize("jq '.age' data.json")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different arg values collide", func(t *testing.T) {
		a := shellshape.Normalize("jq --arg key hello '.' data.json")
		b := shellshape.Normalize("jq --arg user world '.' data.json")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("jq '.name' file.json")
		subshell := shellshape.Normalize("jq $(dangerous-command) file.json")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
