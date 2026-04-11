package shellshape

import "testing"

func TestCal(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare cal", "cal", "cal"},
		{"bare ncal", "ncal", "ncal"},
		{"year only", "cal 2024", "cal N"},
		{"month and year", "cal 3 2024", "cal N N"},
		{"ncal year", "ncal 2024", "ncal N"},

		// Boolean flags
		{"three months", "cal -3", "cal -3"},
		{"yearly", "cal -y", "cal -y"},
		{"julian days", "cal -j", "cal -j"},
		{"no highlight", "cal -h", "cal -h"},
		{"week numbers", "ncal -w", "ncal -w"},
		{"multiple bool flags", "cal -j -y", "cal -j -y"},

		// Flags with arguments
		{"after months", "cal -A 3", "cal -A N"},
		{"before months", "cal -B 2", "cal -B N"},
		{"after and before", "cal -A 2 -B 1", "cal -A N -B N"},
		{"month flag", "cal -m 8", "cal -m N"},
		{"month flag with year", "cal -m 8 2024", "cal -m N N"},
		{"debug date", "ncal -d 2024-06", "ncal -d <date>"},
		{"highlight date", "ncal -H 2024-06-15", "ncal -H <date>"},
		{"country code", "ncal -s FR", "ncal -s FR"},

		// Combined flags and positionals
		{"flags before year", "cal -j 2024", "cal -j N"},
		{"flags with month year", "cal -3 3 2024", "cal -3 N N"},
		{"ncal complex", "ncal -w -A 2 2024", "ncal -w -A N N"},

		// Redirect
		{"with redirect", "cal 2024 > /tmp/cal.txt", "cal N > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different years collide", func(t *testing.T) {
		a := Normalize("cal 2024")
		b := Normalize("cal 1999")
		c := Normalize("cal 2030")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different month-year pairs collide", func(t *testing.T) {
		a := Normalize("cal 1 2024")
		b := Normalize("cal 12 1999")
		if a != b {
			t.Errorf("expected same: %q, %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("cal 2024")
		subshell := Normalize("cal $(date +%Y)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
