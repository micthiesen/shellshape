package shellshape

import "testing"

func TestSu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args (become root)", "su", "su"},
		{"switch user", "su alice", "su <val>"},
		{"login shell", "su - alice", "su - <val>"},
		{"login long flag", "su --login alice", "su --login <val>"},
		{"login short flag", "su -l bob", "su -l <val>"},

		// Flags with arguments
		{"command flag", "su -c 'whoami' alice", "su -c <code> <val>"},
		{"shell flag", "su -s /bin/zsh alice", "su -s <path> <val>"},
		{"group flag", "su -g staff alice", "su -g <val>+"},
		{"long command", "su --command='ls -la' root", "su --command=<val> <val>"},

		// Combined flags
		{"login with command", "su -l -c 'id' bob", "su -l -c <code> <val>"},
		{"preserve env", "su -m alice", "su -m <val>"},
		{"whitelist env", "su -w HOME,PATH alice", "su -w <val>+"},
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
		a := Normalize("su alice")
		b := Normalize("su bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different commands collide", func(t *testing.T) {
		a := Normalize("su -c 'whoami' root")
		b := Normalize("su -c 'id' root")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("su literal-arg")
		subshell := Normalize("su $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
