package shellshape

import "testing"

func TestIconv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"from and to encoding", "iconv -f ISO-8859-1 -t UTF-8 input.txt", "iconv -f <encoding> -t <encoding> <path>"},
		{"from encoding only", "iconv -f SHIFT_JIS input.txt", "iconv -f <encoding> <path>"},
		{"to encoding only", "iconv -t ASCII file.txt", "iconv -t <encoding> <path>"},
		{"list encodings", "iconv -l", "iconv -l"},
		{"output flag", "iconv -f ISO-8859-1 -t UTF-8 input.txt -o output.txt", "iconv -f <encoding> -t <encoding> <path> -o <path>"},

		// Long form flags
		{"long from-code", "iconv --from-code=UTF-16 --to-code=UTF-8 data.txt", "iconv --from-code=<val> --to-code=<val> <path>"},
		{"long output", "iconv -f UTF-8 -t ASCII --output result.txt input.txt", "iconv -f <encoding> -t <encoding> --output <path>+"},

		// Multiple input files
		{"multiple files", "iconv -f ISO-8859-1 -t UTF-8 a.txt b.txt c.txt", "iconv -f <encoding> -t <encoding> <path>+"},

		// Flags with no arguments
		{"discard unconvertible", "iconv -c -f UTF-8 -t ASCII input.txt", "iconv -c -f <encoding> -t <encoding> <path>"},
		{"silent flag", "iconv -s -f UTF-8 -t ASCII file.txt", "iconv -s -f <encoding> -t <encoding> <path>"},

		// Redirect
		{"with redirect", "iconv -f UTF-8 -t ASCII < input.txt > output.txt", "iconv -f <encoding> -t <encoding> < <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different encodings collide", func(t *testing.T) {
		a := Normalize("iconv -f ISO-8859-1 -t UTF-8 input.txt")
		b := Normalize("iconv -f SHIFT_JIS -t ASCII input.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("iconv -f UTF-8 -t ASCII data.csv")
		b := Normalize("iconv -f UTF-8 -t ASCII report.json")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("iconv -f UTF-8 literal.txt")
		subshell := Normalize("iconv -f UTF-8 $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
