package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPkill(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"pkill basic", "pkill nginx", "pkill <pattern>"},
		{"pgrep basic", "pgrep sshd", "pgrep <pattern>"},
		// Match flags (boolean)
		{"full command match", "pkill -f lm-studio", "pkill -f <pattern>"},
		{"exact match", "pgrep -x sshd", "pgrep -x <pattern>"},
		{"newest", "pkill -n firefox", "pkill -n <pattern>"},
		{"oldest", "pkill -o firefox", "pkill -o <pattern>"},
		// Signal flags
		{"numeric signal", "pkill -9 nginx", "pkill -9 <pattern>"},
		{"named signal", "pkill -HUP nginx", "pkill -HUP <pattern>"},
		{"SIG-prefixed signal", "pkill -SIGTERM nginx", "pkill -SIGTERM <pattern>"},
		{"long signal flag", "pkill --signal KILL nginx", "pkill --signal KILL <pattern>"},
		// Boolean flags
		{"count", "pgrep -c nginx", "pgrep -c <pattern>"},
		{"list", "pgrep -l nginx", "pgrep -l <pattern>"},
		{"list full", "pgrep -a nginx", "pgrep -a <pattern>"},
		{"case insensitive", "pgrep -i nginx", "pgrep -i <pattern>"},
		{"negate", "pgrep -v nginx", "pgrep -v <pattern>"},
		// Value-consuming flags
		{"user filter", "pkill -u root nginx", "pkill -u <user> <pattern>"},
		{"group filter", "pkill -g staff nginx", "pkill -g <val> <pattern>"},
		{"gid filter", "pkill -G 1000 nginx", "pkill -G <val> <pattern>"},
		{"parent filter", "pkill -P 1234 nginx", "pkill -P <val> <pattern>"},
		{"session filter", "pkill -s 0 nginx", "pkill -s <val> <pattern>"},
		{"terminal filter", "pkill -t tty1 nginx", "pkill -t <val> <pattern>"},
		// Combined flags
		{"signal and match", "pkill -9 -f lm-studio", "pkill -9 -f <pattern>"},
		{"multiple flags", "pgrep -l -u root -f nginx", "pgrep -l -u <user> -f <pattern>"},
		// Subshell preserved
		{"subshell standalone", "pkill $(echo nginx)", "pkill $(echo <str>)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different patterns → same shape
	t.Run("different patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("pkill -f lm-studio")
		b := shellshape.Normalize("pkill -f nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different users collide", func(t *testing.T) {
		a := shellshape.Normalize("pkill -u root nginx")
		b := shellshape.Normalize("pkill -u www-data apache2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pkill nginx")
		subshell := shellshape.Normalize("pkill $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
