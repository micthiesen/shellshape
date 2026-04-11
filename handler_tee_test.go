package shellshape

import "testing"

func TestTee(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "tee", "tee"},
		{"single file", "tee output.txt", "tee <path>"},
		{"absolute path", "tee /tmp/log.txt", "tee <path>"},
		{"tilde path", "tee ~/logs/out.txt", "tee <path>"},

		// Flags
		{"append flag", "tee -a output.txt", "tee -a <path>"},
		{"ignore sigint", "tee -i output.txt", "tee -i <path>"},
		{"bundled flags", "tee -ai output.txt", "tee -ai <path>"},

		// Multiple positionals
		{"two files", "tee file1.txt file2.txt", "tee <path>+"},
		{"three files with flag", "tee -a file1.txt file2.txt file3.txt", "tee -a <path>+"},

		// Bare word (no extension) still becomes <path>
		{"bare word file", "tee mylog", "tee <path>"},
		{"bare words multiple", "tee a b c", "tee <path>+"},

		// With redirect
		{"with redirect", "tee output.txt > /dev/null", "tee <path> > <path>"},

		// Piped usage shape (just the tee part)
		{"flag only no file", "tee -a", "tee -a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different file paths produce the same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("tee -a /var/log/app.log")
		b := Normalize("tee -a /tmp/debug.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Safety test: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tee literal-arg")
		subshell := Normalize("tee $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
