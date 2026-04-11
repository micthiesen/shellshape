package shellshape

import "testing"

func TestKill(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single pid", "kill 1234", "kill <pid>"},
		{"multiple pids", "kill 1234 5678 9012", "kill <pid>+"},
		// Signal flags
		{"numeric signal", "kill -9 1234", "kill -9 <pid>"},
		{"named signal", "kill -HUP 1234", "kill -HUP <pid>"},
		{"SIG-prefixed signal", "kill -SIGKILL 1234", "kill -SIGKILL <pid>"},
		{"dash-s flag", "kill -s TERM 1234", "kill -s TERM <pid>"},
		{"dash-s with multiple pids", "kill -s HUP 1234 5678", "kill -s HUP <pid>+"},
		// List mode
		{"list signals", "kill -l", "kill -l"},
		{"list with exit status", "kill -l 1", "kill -l N"},
		// Special pids
		{"job spec", "kill %1", "kill <pid>"},
		{"double dash", "kill -- -1", "kill -- <pid>"},
		{"pid zero", "kill 0", "kill <pid>"},
		{"negative pid", "kill -- -117", "kill -- <pid>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different PIDs → same shape
	t.Run("different pids collide", func(t *testing.T) {
		a := Normalize("kill -9 1234")
		b := Normalize("kill -9 5678")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pid counts collide", func(t *testing.T) {
		a := Normalize("kill 1234 5678")
		b := Normalize("kill 111 222 333")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("kill 1234")
		subshell := Normalize("kill $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
