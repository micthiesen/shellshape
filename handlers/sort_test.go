package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSort(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "sort file.txt", "sort <path>"},
		{"multiple files", "sort a.txt b.txt c.txt", "sort <path>+"},
		{"no args", "sort", "sort"},
		{"stdin (no files)", "sort -n", "sort -n"},

		// Boolean flags
		{"numeric reverse", "sort -n -r data.csv", "sort -n -r <path>"},
		{"bundled flags", "sort -nr data.csv", "sort -nr <path>"},
		{"unique", "sort -u file.txt", "sort -u <path>"},
		{"check sorted", "sort -c file.txt", "sort -c <path>"},
		{"human numeric", "sort -h file.txt", "sort -h <path>"},

		// Flags with arguments
		{"key spec", "sort -k 2,3 file.txt", "sort -k <key> <path>"},
		{"multiple keys", "sort -k 1,1 -k 2,2nr file.txt", "sort -k <key> -k <key> <path>"},
		{"field separator", "sort -t , file.txt", "sort -t <sep> <path>"},
		{"field separator colon", "sort -t : /etc/passwd", "sort -t <sep> <path>"},
		{"output file", "sort -o output.txt input.txt", "sort -o <path>+"},
		{"buffer size", "sort -S 1G file.txt", "sort -S <size> <path>"},
		{"temp dir", "sort -T /tmp file.txt", "sort -T <path>+"},

		// Long flags with =
		{"long key", "sort --key=2,3 file.txt", "sort --key=<val> <path>"},
		{"long output", "sort --output=result.txt file.txt", "sort --output=<val> <path>"},
		{"long separator", "sort --field-separator=, file.txt", "sort --field-separator=<val> <path>"},

		// Long flags without =
		{"long key no eq", "sort --key 2,3 file.txt", "sort --key <key> <path>"},
		{"long output no eq", "sort --output result.txt file.txt", "sort --output <path>+"},

		// Combined usage
		{"complex", "sort -t : -k 3,3n -k 1,1 /etc/passwd", "sort -t <sep> -k <key> -k <key> <path>"},

		// Redirects
		{"redirect out", "sort file.txt > sorted.txt", "sort <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("sort -k 2,3 data.csv")
		b := shellshape.Normalize("sort -k 2,3 results.tsv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different keys collide", func(t *testing.T) {
		a := shellshape.Normalize("sort -k 1,1 file.txt")
		b := shellshape.Normalize("sort -k 3,5nr file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sort literal-arg")
		subshell := shellshape.Normalize("sort $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
