package shellshape

import "testing"

func TestLess(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare less", "less", "less"},
		{"single file", "less file.txt", "less <path>"},
		{"absolute path", "less /var/log/syslog", "less <path>"},
		{"relative path", "less ./README.md", "less <path>"},

		// Boolean flags
		{"line numbers", "less -N file.txt", "less -N <path>"},
		{"no wrap", "less -S file.txt", "less -S <path>"},
		{"raw control", "less -R file.txt", "less -R <path>"},
		{"no termcap init", "less -X file.txt", "less -X <path>"},
		{"combined flags", "less -RXS file.txt", "less -RXS <path>"},
		{"follow mode", "less +F file.txt", "less +F <path>"},

		// Flags with numeric arguments
		{"-x tab stops", "less -x 4 file.txt", "less -x N <path>"},
		{"-j target line", "less -j 10 file.txt", "less -j N <path>"},
		{"-z window size", "less -z 25 file.txt", "less -z N <path>"},

		// Flags with value arguments
		{"-p pattern", "less -p 'error' file.txt", "less -p <val> <path>"},
		{"-t tag", "less -t my_function file.txt", "less -t <val> <path>"},

		// Flags with path arguments
		{"-k lesskey", "less -k ~/.lesskey file.txt", "less -k <path>+"},
		{"-o log file", "less -o /tmp/less.log file.txt", "less -o <path>+"},

		// Multiple files
		{"two files", "less file1.txt file2.txt", "less <path>+"},
		{"three files", "less a.go b.go c.go", "less <path>+"},

		// With redirect
		{"redirect", "less file.txt > /tmp/out", "less <path> > <path>"},

		// more alias
		{"more basic", "more file.txt", "more <path>"},
		{"more flags", "more -d file.txt", "more -d <path>"},
		{"more multiple", "more a.txt b.txt", "more <path>+"},
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
		a := Normalize("less /etc/hosts")
		b := Normalize("less /var/log/syslog")
		c := Normalize("less ~/.bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different patterns collide", func(t *testing.T) {
		a := Normalize("less -p 'error' file.txt")
		b := Normalize("less -p 'warning' file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("less literal-arg")
		subshell := Normalize("less $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
