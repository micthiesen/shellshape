package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare date", "date", "date"},
		{"with format", "date +%Y-%m-%d", "date <fmt>"},
		{"format with time", "date +%H:%M:%S", "date <fmt>"},
		{"utc flag", "date -u", "date -u"},
		{"rfc email", "date -R", "date -R"},

		// Flags with arguments
		{"date string", `date -d "2024-01-01"`, "date -d <date-str>"},
		{"date string with format", `date -d yesterday +%Y-%m-%d`, "date -d <date-str> <fmt>"},
		{"set date", `date -s "2024-01-01 12:00"`, "date -s <date-str>"},
		{"reference file", `date -r /etc/passwd`, "date -r <date-str>"},
		{"file flag", `date -f /tmp/dates.txt`, "date -f <path>"},

		// Combined flags
		{"utc with format", "date -u +%s", "date -u <fmt>"},
		{"iso flag", "date -I", "date -I"},
		{"rfc3339 equals", "date --rfc-3339=seconds", "date --rfc-3339=seconds"},
		{"long date flag", `date --date="last friday"`, "date --date=<date-str>"},
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
	t.Run("different formats collide", func(t *testing.T) {
		a := shellshape.Normalize("date +%Y-%m-%d")
		b := shellshape.Normalize("date +%H:%M:%S")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different date strings collide", func(t *testing.T) {
		a := shellshape.Normalize(`date -d "2024-01-01"`)
		b := shellshape.Normalize(`date -d "last friday"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("date +%Y-%m-%d")
		subshell := shellshape.Normalize("date $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
