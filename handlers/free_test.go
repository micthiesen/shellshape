package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFree(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "free", "free"},
		{"human readable", "free -h", "free -h"},
		{"megabytes", "free -m", "free -m"},
		{"gigabytes", "free -g", "free -g"},
		{"bytes", "free -b", "free -b"},
		{"wide output", "free -w", "free -w"},
		{"total line", "free -t", "free -t"},
		{"lohi stats", "free -l", "free -l"},

		// Combined boolean flags
		{"combined flags", "free -ht", "free -ht"},
		{"combined flags wide", "free -htw", "free -htw"},
		{"giga total", "free -gt", "free -gt"},

		// Long boolean flags
		{"long human", "free --human", "free --human"},
		{"long wide", "free --wide", "free --wide"},
		{"long total", "free --total", "free --total"},

		// Flags with arguments (numeric)
		{"seconds short", "free -s 2", "free -s N"},
		{"seconds long", "free --seconds 10", "free --seconds N"},
		{"count short", "free -c 5", "free -c N"},
		{"count long", "free --count 3", "free --count N"},

		// Combined flags ending in consuming flag
		{"human with seconds", "free -hs 5", "free -hs N"},

		// Multiple flags with arguments
		{"seconds and count", "free -s 2 -c 10", "free -s N -c N"},
		{"human seconds count", "free -h -s 2 -c 3", "free -h -s N -c N"},

		// Long flag with equals
		{"seconds equals", "free --seconds=5", "free --seconds=<val>"},
		{"count equals", "free --count=10", "free --count=<val>"},

		// Consuming flag at end of input (no value follows)
		{"seconds short trailing", "free -s", "free -s"},
		{"count short trailing", "free -c", "free -c"},
		{"seconds long trailing", "free --seconds", "free --seconds"},
		{"count long trailing", "free --count", "free --count"},

		// Combined flag ending in consuming letter at end of input
		{"combined trailing s", "free -hs", "free -hs"},
		{"combined trailing c", "free -hc", "free -hc"},

		// Combined flag ending in consuming letter with subshell value
		{"combined s subshell", "free -hs $(calc)", "free -hs $(calc)"},
		{"combined c subshell", "free -hc $(calc)", "free -hc $(calc)"},

		// Subshell as standalone token
		{"subshell standalone", "free $(flags)", "free $(flags)"},

		// Consuming flag with subshell value
		{"count subshell", "free -c $(echo 5)", "free -c $(echo <str>)"},

		// Edge cases
		{"redirect", "free -h > output.txt", "free -h > <path>"},
		{"pipe ignored in shape", "free -m", "free -m"},
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
	t.Run("different intervals collide", func(t *testing.T) {
		a := shellshape.Normalize("free -s 2")
		b := shellshape.Normalize("free -s 30")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different counts collide", func(t *testing.T) {
		a := shellshape.Normalize("free -c 5")
		b := shellshape.Normalize("free -c 100")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("free -s 5")
		subshell := shellshape.Normalize("free -s $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
