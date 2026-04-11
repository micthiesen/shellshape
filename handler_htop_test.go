package shellshape

import "testing"

func TestHtop(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "htop", "htop"},
		{"tree mode", "htop -t", "htop -t"},
		{"help", "htop --help", "htop --help"},
		{"version", "htop -V", "htop -V"},

		// Flags with numeric arguments
		{"delay short", "htop -d 10", "htop -d N"},
		{"delay long", "htop --delay 5", "htop --delay N"},
		{"highlight changes", "htop -H 20", "htop -H N"},
		{"highlight changes long", "htop --highlight-changes 15", "htop --highlight-changes N"},

		// Flags with value arguments
		{"user short", "htop -u root", "htop -u <user>"},
		{"user long", "htop --user nobody", "htop --user <user>"},
		{"pid short", "htop -p 1234,5678", "htop -p <pid>"},
		{"pid long", "htop --pid 42", "htop --pid <pid>"},
		{"sort short", "htop -s PERCENT_CPU", "htop -s <col>"},
		{"sort long", "htop --sort-key PERCENT_MEM", "htop --sort-key <col>"},
		{"filter short", "htop -F firefox", "htop -F <filter>"},
		{"filter long", "htop --filter nginx", "htop --filter <filter>"},

		// Boolean flags
		{"no color", "htop -C", "htop -C"},
		{"no color long", "htop --no-color", "htop --no-color"},
		{"no mouse", "htop -M", "htop -M"},
		{"no unicode", "htop -U", "htop -U"},
		{"readonly", "htop --readonly", "htop --readonly"},

		// Combined flags
		{"tree with user", "htop -t -u www-data", "htop -t -u <user>"},
		{"delay and sort", "htop -d 20 -s PERCENT_CPU", "htop -d N -s <col>"},
		{"tree no color user", "htop -t -C -u admin", "htop -t -C -u <user>"},

		// Redirect
		{"redirect", "htop -p 123 > out.txt", "htop -p <pid> > <path>"},
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
	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("htop -u root")
		b := Normalize("htop -u nobody")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pids collide", func(t *testing.T) {
		a := Normalize("htop -p 1234")
		b := Normalize("htop -p 9999,8888")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different delays collide", func(t *testing.T) {
		a := Normalize("htop -d 5")
		b := Normalize("htop -d 50")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("htop -u root")
		subshell := Normalize("htop -u $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
