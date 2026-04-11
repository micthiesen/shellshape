package shellshape

import "testing"

func TestNl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "nl", "nl"},
		{"single file", "nl file.txt", "nl <path>"},
		{"stdin dash", "nl -", "nl -"},
		{"multiple files", "nl file1.txt file2.txt", "nl <path>+"},

		// Boolean flag
		{"persist flag", "nl -p file.txt", "nl -p <path>"},

		// Flags with value arguments
		{"body style all", "nl -b a file.txt", "nl -b a <path>"},
		{"body style none", "nl -b n file.txt", "nl -b n <path>"},
		{"body style pattern", "nl -b p'foo.*bar' file.txt", "nl -b p<pattern> <path>"},
		{"footer style", "nl -f t file.txt", "nl -f t <path>"},
		{"header style", "nl -h a file.txt", "nl -h a <path>"},
		{"delimiter", "nl -d '::' file.txt", "nl -d <str> <path>"},
		{"number format", "nl -n rz file.txt", "nl -n rz <path>"},
		{"separator", "nl -s ': ' file.txt", "nl -s <str> <path>"},

		// Numeric flags
		{"increment", "nl -i 5 file.txt", "nl -i N <path>"},
		{"join blank", "nl -l 2 file.txt", "nl -l N <path>"},
		{"start number", "nl -v 10 file.txt", "nl -v N <path>"},
		{"width", "nl -w 8 file.txt", "nl -w N <path>"},

		// Combined flags
		{"combined", "nl -ba -nrz -s '-> ' -v 0 -w 4 file.txt",
			"nl -b a -n rz -s <str> -v N -w N <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("nl -ba /var/log/syslog")
		b := Normalize("nl -ba /tmp/output.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different patterns collide", func(t *testing.T) {
		a := Normalize("nl -b p'error' file.txt")
		b := Normalize("nl -b p'warning' file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numbers collide", func(t *testing.T) {
		a := Normalize("nl -v 1 -i 2 file.txt")
		b := Normalize("nl -v 100 -i 10 file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("nl literal-arg")
		subshell := Normalize("nl $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
