package shellshape

import "testing"

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

		// Edge cases
		{"redirect", "free -h > output.txt", "free -h > <path>"},
		{"pipe ignored in shape", "free -m", "free -m"},
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
	t.Run("different intervals collide", func(t *testing.T) {
		a := Normalize("free -s 2")
		b := Normalize("free -s 30")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different counts collide", func(t *testing.T) {
		a := Normalize("free -c 5")
		b := Normalize("free -c 100")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("free -s 5")
		subshell := Normalize("free -s $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
