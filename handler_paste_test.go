package shellshape

import "testing"

func TestPaste(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare paste", "paste", "paste"},
		{"two files", "paste file1.txt file2.txt", "paste <path>+"},
		{"three files", "paste a.txt b.txt c.txt", "paste <path>+"},
		{"absolute paths", "paste /tmp/a.txt /tmp/b.txt", "paste <path>+"},

		// Stdin marker
		{"stdin columns", "paste - - -", "paste - - -"},
		{"stdin with file", "paste - file.txt", "paste - <path>"},

		// Serial flag
		{"serial flag", "paste -s file.txt", "paste -s <path>"},
		{"serial with multiple", "paste -s file1.txt file2.txt", "paste -s <path>+"},

		// Delimiter flag
		{"delimiter separate", "paste -d : file1.txt file2.txt", "paste -d <delim> <path>+"},
		{"delimiter tab", `paste -d '\t' file1.txt file2.txt`, `paste -d <delim> <path>+`},
		{"delimiter fused", "paste -d, file1.txt file2.txt", "paste -d <delim> <path>+"},

		// Combined flags
		{"serial and delimiter", "paste -s -d : file.txt", "paste -s -d <delim> <path>"},
		{"serial and delimiter stdin", "paste -s -d : -", "paste -s -d <delim> -"},

		// Long flags (GNU coreutils)
		{"long serial", "paste --serial file.txt", "paste --serial <path>"},
		{"long delimiters", "paste --delimiters=: file1.txt file2.txt", "paste --delimiters=<val> <path>+"},

		// Redirects
		{"with redirect", "paste file1.txt file2.txt > out.txt", "paste <path>+ > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("paste /tmp/a.txt /tmp/b.txt")
		b := Normalize("paste /var/x.csv /var/y.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different delimiters collide", func(t *testing.T) {
		a := Normalize("paste -d : file1.txt file2.txt")
		b := Normalize("paste -d , file1.txt file2.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("paste literal-arg")
		subshell := Normalize("paste $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
