package shellshape

import "testing"

func TestTac(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare tac", "tac", "tac"},
		{"single file", "tac file.txt", "tac <path>"},
		{"absolute path", "tac /var/log/syslog", "tac <path>"},
		{"relative path", "tac ./data.csv", "tac <path>"},

		// Boolean flags
		{"before flag", "tac -b file.txt", "tac -b <path>"},
		{"regex flag", "tac -r file.txt", "tac -r <path>"},
		{"both boolean flags", "tac -b -r file.txt", "tac -b -r <path>"},

		// Separator flag (consumes next token)
		{"separator", "tac -s '---' file.txt", "tac -s <sep> <path>"},
		{"long separator", "tac --separator='\\n\\n' file.txt", "tac --separator=<sep> <path>"},
		{"separator with regex", "tac -r -s '^Section' file.txt", "tac -r -s <sep> <path>"},

		// Multiple positionals
		{"two files", "tac file1.txt file2.txt", "tac <path>+"},
		{"three files", "tac a.log b.log c.log", "tac <path>+"},

		// Redirects
		{"with redirect", "tac file.txt > /tmp/reversed", "tac <path> > <path>"},

		// Pipeline
		{"in pipeline", "tac file.txt && cat file.txt", "tac <path> && cat <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different file paths → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("tac /etc/hosts")
		b := Normalize("tac /var/log/syslog")
		c := Normalize("tac ~/.bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// COLLISION TEST: different separators → same shape
	t.Run("different separators collide", func(t *testing.T) {
		a := Normalize("tac -s '---' file.txt")
		b := Normalize("tac -s '===' file.txt")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tac literal-arg")
		subshell := Normalize("tac $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
