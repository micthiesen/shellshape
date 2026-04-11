package shellshape

import "testing"

func TestCat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare cat", "cat", "cat"},
		{"single file", "cat file.txt", "cat <path>"},
		{"absolute path", "cat /etc/hosts", "cat <path>"},
		{"relative path", "cat ./config.yaml", "cat <path>"},

		// Flags (all boolean)
		{"number lines", "cat -n file.txt", "cat -n <path>"},
		{"number non-blank", "cat -b file.txt", "cat -b <path>"},
		{"squeeze blank", "cat -s file.txt", "cat -s <path>"},
		{"show ends", "cat -e file.txt", "cat -e <path>"},
		{"show tabs", "cat -t file.txt", "cat -t <path>"},
		{"combined flags", "cat -ns file.txt", "cat -ns <path>"},

		// Multiple positionals (collapsed by collapseRepeatedPlaceholders)
		{"two files", "cat file1.txt file2.txt", "cat <path>+"},
		{"three files", "cat a.go b.go c.go", "cat <path>+"},

		// Stdin marker
		{"stdin dash", "cat -", "cat -"},

		// Redirects
		{"with redirect", "cat file.txt > /tmp/out", "cat <path> > <path>"},

		// Pipeline
		{"in pipeline", "cat file.txt && cat other.txt", "cat <path> && cat <path>"},
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
		a := Normalize("cat /etc/hosts")
		b := Normalize("cat /var/log/syslog")
		c := Normalize("cat ~/.bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("cat literal-arg")
		subshell := Normalize("cat $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestRev(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare rev", "rev", "rev"},
		{"single file", "rev file.txt", "rev <path>"},
		{"absolute path", "rev /etc/hosts", "rev <path>"},
		{"relative path", "rev ./data.csv", "rev <path>"},

		// Multiple positionals
		{"two files", "rev file1.txt file2.txt", "rev <path>+"},
		{"three files", "rev a.txt b.txt c.txt", "rev <path>+"},

		// Stdin marker
		{"stdin dash", "rev -", "rev -"},

		// Redirects
		{"with redirect", "rev file.txt > /tmp/out", "rev <path> > <path>"},

		// Pipeline
		{"in pipeline", "rev file.txt && rev other.txt", "rev <path> && rev <path>"},
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
		a := Normalize("rev /etc/hosts")
		b := Normalize("rev /var/log/syslog")
		c := Normalize("rev ~/.bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("rev literal-arg")
		subshell := Normalize("rev $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
