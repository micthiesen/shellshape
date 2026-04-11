package shellshape

import "testing"

func TestCut(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"field selection", "cut -f 1 file.txt", "cut -f <range> <path>"},
		{"byte selection", "cut -b 1-5 file.txt", "cut -b <range> <path>"},
		{"char selection", "cut -c 1-16,26-38", "cut -c <range>"},
		{"delimiter and field", "cut -d : -f 1,7 /etc/passwd", "cut -d <delim> -f <range> <path>"},
		{"suppress flag", "cut -d : -f 1 -s file.txt", "cut -d <delim> -f <range> -s <path>"},

		// Fused flag-value forms
		{"fused delimiter", "cut -d: -f 1,7 /etc/passwd", "cut -d <delim> -f <range> <path>"},
		{"fused field", "cut -f1,3 data.csv", "cut -f <range> <path>"},
		{"fused byte range", "cut -b1-5 file.txt", "cut -b <range> <path>"},
		{"fused char range", "cut -c1-16,26-38", "cut -c <range>"},

		// Long form flags (GNU cut)
		{"long delimiter", "cut --delimiter=: --fields=1,7 /etc/passwd", "cut --delimiter=<val> --fields=<val> <path>"},

		// Multiple files
		{"multiple files", "cut -f 1 file1.txt file2.txt", "cut -f <range> <path>+"},

		// Whitespace delimiter
		{"whitespace delimiter", "cut -w -f 1 file.txt", "cut -w -f <range> <path>"},

		// Boolean flag -n
		{"no split multibyte", "cut -b 1-5 -n file.txt", "cut -b <range> -n <path>"},

		// Redirect
		{"with redirect", "cut -f 1 > output.txt", "cut -f <range> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different delimiters collide", func(t *testing.T) {
		a := Normalize("cut -d : -f 1,7 /etc/passwd")
		b := Normalize("cut -d , -f 1,7 /etc/passwd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ranges collide", func(t *testing.T) {
		a := Normalize("cut -f 1,3 file.txt")
		b := Normalize("cut -f 2,5,7 file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("cut -f 1 file.txt")
		subshell := Normalize("cut -f 1 $(echo file.txt)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
