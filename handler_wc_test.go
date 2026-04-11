package shellshape

import "testing"

func TestWc(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "wc file.txt", "wc <path>"},
		{"no args stdin", "wc", "wc"},
		{"line count", "wc -l file.txt", "wc -l <path>"},
		{"word count", "wc -w file.txt", "wc -w <path>"},
		{"byte count", "wc -c file.txt", "wc -c <path>"},
		{"char count", "wc -m file.txt", "wc -m <path>"},
		{"longest line", "wc -L file.txt", "wc -L <path>"},
		// Bundled flags
		{"bundled flags", "wc -lw file.txt", "wc -lw <path>"},
		{"all flags bundled", "wc -mlw report1 report2", "wc -mlw <path>+"},
		// Multiple files
		{"multiple files", "wc -l report1 report2 report3", "wc -l <path>+"},
		{"multiple files no flags", "wc foo.txt bar.txt baz.txt", "wc <path>+"},
		// Long flags
		{"long lines flag", "wc --lines myfile.log", "wc --lines <path>"},
		// Flag only (stdin)
		{"flag only stdin", "wc -l", "wc -l"},
		// Redirect
		{"redirect", "wc -l < input.txt", "wc -l < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different file names → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("wc -l /var/log/syslog")
		b := Normalize("wc -l /tmp/output.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different multi files collide", func(t *testing.T) {
		a := Normalize("wc -w report1.txt report2.txt")
		b := Normalize("wc -w data.csv results.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("wc file.txt")
		subshell := Normalize("wc $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
