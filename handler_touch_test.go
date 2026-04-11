package shellshape

import "testing"

func TestTouch(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "touch file.txt", "touch <path>"},
		{"multiple files", "touch foo.txt bar.txt baz.txt", "touch <path>+"},
		{"absolute path", "touch /tmp/newfile", "touch <path>"},

		// Boolean flags
		{"no-create", "touch -c file.txt", "touch -c <path>"},
		{"access time only", "touch -a file.txt", "touch -a <path>"},
		{"modification time only", "touch -m file.txt", "touch -m <path>"},
		{"symlink flag", "touch -h file.txt", "touch -h <path>"},
		{"bundled boolean flags", "touch -acm file.txt", "touch -acm <path>"},

		// Flags with arguments
		{"timestamp -t", "touch -t 202301011200.00 file.txt", "touch -t <val> <path>"},
		{"date -d", "touch -d 2023-01-01T12:00:00 file.txt", "touch -d <val> <path>"},
		{"reference -r", "touch -r reference.txt target.txt", "touch -r <path>+"},
		{"adjust -A", "touch -A 0100 file.txt", "touch -A <val> <path>"},

		// Edge cases
		{"redirect", "touch file.txt 2>/dev/null", "touch <path> 2>/dev/null"},
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
		a := Normalize("touch /home/user/foo.txt")
		b := Normalize("touch /var/log/bar.log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different timestamps collide", func(t *testing.T) {
		a := Normalize("touch -t 202301011200.00 file.txt")
		b := Normalize("touch -t 199912312359.59 other.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("touch literal-arg")
		subshell := Normalize("touch $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
