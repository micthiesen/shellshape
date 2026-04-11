package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUniq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "uniq", "uniq"},
		{"input file", "uniq file.txt", "uniq <path>"},
		{"input and output files", "uniq input.txt output.txt", "uniq <path>+"},

		// Boolean flags
		{"count", "uniq -c file.txt", "uniq -c <path>"},
		{"repeated", "uniq -d file.txt", "uniq -d <path>"},
		{"unique", "uniq -u file.txt", "uniq -u <path>"},
		{"ignore case", "uniq -i file.txt", "uniq -i <path>"},
		{"combined flags", "uniq -d -i file.txt", "uniq -d -i <path>"},

		// Flags with numeric arguments
		{"skip fields", "uniq -f 3 file.txt", "uniq -f N <path>"},
		{"skip chars", "uniq -s 5 file.txt", "uniq -s N <path>"},
		{"skip fields long", "uniq --skip-fields 2 data.csv", "uniq --skip-fields N <path>"},
		{"skip chars long", "uniq --skip-chars 10 data.csv", "uniq --skip-chars N <path>"},
		{"skip fields and count", "uniq -f 3 -c file.txt", "uniq -f N -c <path>"},

		// -D with septype (structural keyword, kept verbatim)
		{"all-repeated", "uniq -D file.txt", "uniq -D <path>"},
		{"all-repeated prepend", "uniq -D prepend file.txt", "uniq -D prepend <path>"},
		{"all-repeated separate", "uniq -D separate file.txt", "uniq -D separate <path>"},
		{"all-repeated none", "uniq -D none file.txt", "uniq -D none <path>"},

		// Redirect
		{"with redirect", "uniq -c file.txt > out.txt", "uniq -c <path> > <path>"},
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
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("uniq -c /var/log/syslog")
		b := shellshape.Normalize("uniq -c /tmp/data.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different skip numbers collide", func(t *testing.T) {
		a := shellshape.Normalize("uniq -f 3 file.txt")
		b := shellshape.Normalize("uniq -f 10 other.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("uniq literal-arg")
		subshell := shellshape.Normalize("uniq $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
