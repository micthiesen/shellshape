package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"translate lower to upper", "tr a-z A-Z", "tr <set>+"},
		{"translate with classes", "tr '[:lower:]' '[:upper:]'", "tr <set>+"},
		{"delete characters", "tr -d '\\n'", "tr -d <set>"},
		{"squeeze spaces", "tr -s ' '", "tr -s <set>"},
		{"complement and delete", "tr -cd '[:print:]'", "tr -cd <set>"},
		{"complement squeeze translate", "tr -cs '[:alpha:]' '\\n'", "tr -cs <set>+"},
		{"delete and squeeze", "tr -ds '[:lower:]' '[:upper:]'", "tr -ds <set>+"},
		{"unbuffered flag", "tr -u a-z A-Z", "tr -u <set>+"},

		// Flags and positionals mixed
		{"separate flags", "tr -c -s '[:alpha:]' '\\n'", "tr -c -s <set>+"},

		// Redirect
		{"with redirect", "tr a-z A-Z < input.txt", "tr <set>+ < <path>"},
		{"with pipe-style redirect", "tr -d '\\r' < file.txt > out.txt", "tr -d <set> < <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different character sets collide", func(t *testing.T) {
		a := shellshape.Normalize("tr a-z A-Z")
		b := shellshape.Normalize("tr 0-9 a-j")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different delete targets collide", func(t *testing.T) {
		a := shellshape.Normalize("tr -d '\\n'")
		b := shellshape.Normalize("tr -d '\\r'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("tr literal-arg x")
		subshell := shellshape.Normalize("tr $(dangerous-command) x")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
