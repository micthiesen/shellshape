package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKillall(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple", "killall firefox", "killall <pattern>"},
		{"multiple procs", "killall firefox chrome", "killall <pattern>+"},

		// Signal flags (structural, kept verbatim)
		{"signal flag", "killall -9 nginx", "killall -9 <pattern>"},
		{"signal name", "killall -s TERM nginx", "killall -s TERM <pattern>"},
		{"signal long", "killall --signal HUP nginx", "killall --signal HUP <pattern>"},
		{"SIGKILL style", "killall -KILL firefox", "killall -KILL <pattern>"},

		// Boolean flags
		{"exact", "killall -e firefox", "killall -e <pattern>"},
		{"interactive", "killall -i nginx", "killall -i <pattern>"},
		{"regexp", "killall -r 'fire.*'", "killall -r <pattern>"},
		{"wait quiet", "killall -w -q firefox", "killall -w -q <pattern>"},

		// User flag
		{"user", "killall -u root nginx", "killall -u <val> <pattern>"},
		{"user long", "killall --user michael firefox", "killall --user <val> <pattern>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different process names collide", func(t *testing.T) {
		a := shellshape.Normalize("killall firefox")
		b := shellshape.Normalize("killall nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("killall firefox")
		subshell := shellshape.Normalize("killall $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
