package shellshape

import "testing"

func TestUname(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "uname", "uname"},
		{"all info", "uname -a", "uname -a"},
		{"kernel name", "uname -s", "uname -s"},
		{"nodename", "uname -n", "uname -n"},
		{"kernel release", "uname -r", "uname -r"},

		// Combined short flags
		{"combined srm", "uname -srm", "uname -srm"},
		{"combined snrv", "uname -snrv", "uname -snrv"},

		// Long flags
		{"long kernel-name", "uname --kernel-name", "uname --kernel-name"},
		{"long all", "uname --all", "uname --all"},
		{"long machine", "uname --machine", "uname --machine"},
		{"long kernel-release", "uname --kernel-release", "uname --kernel-release"},

		// Redirects
		{"redirect", "uname -a > /tmp/sysinfo.txt", "uname -a > <path>"},
		{"redirect release", "uname -r >> /var/log/info.log", "uname -r >> <path>"},

		// Edge cases: unexpected positional
		{"unexpected positional", "uname foo", "uname <str>"},
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
	t.Run("different positionals collide", func(t *testing.T) {
		a := Normalize("uname foo")
		b := Normalize("uname bar")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("uname foo")
		subshell := Normalize("uname $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
