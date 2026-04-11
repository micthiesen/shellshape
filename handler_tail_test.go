package shellshape

import "testing"

func TestTail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "tail file.txt", "tail <path>"},
		{"absolute path", "tail /var/log/syslog", "tail <path>"},
		{"no args", "tail", "tail"},

		// Flags with numeric arguments
		{"-n lines", "tail -n 20 file.txt", "tail -n N <path>"},
		{"-c bytes", "tail -c 100 file.txt", "tail -c N <path>"},
		{"-b blocks", "tail -b 5 file.txt", "tail -b N <path>"},
		{"-n from beginning", "tail -n +50 file.txt", "tail -n N <path>"},

		// Boolean flags
		{"-f follow", "tail -f /var/log/messages", "tail -f <path>"},
		{"-F follow rename", "tail -F /var/log/messages", "tail -F <path>"},
		{"-r reverse", "tail -r file.txt", "tail -r <path>"},
		{"-q quiet", "tail -q file1.txt file2.txt", "tail -q <path>+"},
		{"-v verbose", "tail -v file.txt", "tail -v <path>"},

		// Combined flags
		{"-f with -n", "tail -f -n 50 /var/log/app.log", "tail -f -n N <path>"},

		// Multiple files
		{"multiple files", "tail file1.txt file2.txt file3.txt", "tail <path>+"},

		// Long flags with =
		{"--lines=N", "tail --lines=20 file.txt", "tail --lines=<val> <path>"},
		{"--bytes=N", "tail --bytes=100 file.txt", "tail --bytes=<val> <path>"},

		// Redirect
		{"redirect", "tail -n 10 file.txt > output.txt", "tail -n N <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("tail -n 20 /var/log/syslog")
		b := Normalize("tail -n 50 /var/log/messages")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tail file.txt")
		subshell := Normalize("tail $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
