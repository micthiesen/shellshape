package shellshape

import "testing"

func TestMkdir(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple bare name", "mkdir mydir", "mkdir mydir"},
		{"simple path", "mkdir /tmp/mydir", "mkdir <path>"},
		{"relative path", "mkdir ./subdir", "mkdir <path>"},
		{"nested path", "mkdir some/nested/dir", "mkdir <path>"},

		// Flags
		{"create parents", "mkdir -p /tmp/a/b/c", "mkdir -p <path>"},
		{"verbose", "mkdir -v /tmp/mydir", "mkdir -v <path>"},
		{"bundled flags", "mkdir -pv /tmp/mydir", "mkdir -pv <path>"},

		// Mode flag with argument
		{"mode numeric", "mkdir -m 755 /tmp/mydir", "mkdir -m <mode> <path>"},
		{"mode symbolic", "mkdir -m u+rwx /tmp/mydir", "mkdir -m <mode> <path>"},
		{"mode with parents", "mkdir -pm 700 /tmp/a/b", "mkdir -pm <mode> <path>"},

		// Multiple positionals
		{"multiple paths", "mkdir -p /tmp/a /tmp/b /tmp/c", "mkdir -p <path>+"},
		{"multiple bare names", "mkdir foo bar baz", "mkdir foo bar baz"},

		// Redirect
		{"with redirect", "mkdir /tmp/mydir 2>/dev/null", "mkdir <path> 2>/dev/null"},
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
	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("mkdir -p /home/alice/project")
		b := Normalize("mkdir -p /var/log/myapp")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different modes collide", func(t *testing.T) {
		a := Normalize("mkdir -m 755 /tmp/a")
		b := Normalize("mkdir -m 700 /tmp/b")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("mkdir literal-arg")
		subshell := Normalize("mkdir $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
