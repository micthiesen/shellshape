package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestScreen(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare screen", "screen", "screen"},
		{"list sessions", "screen -ls", "screen -ls"},
		{"list sessions long", "screen -list", "screen -list"},

		// Named session
		{"new named session", "screen -S myproject", "screen -S <val>"},
		{"detached named session", "screen -dmS background-task ./run.sh", "screen -dmS <val> <path>"},

		// Reattach
		{"reattach any", "screen -r", "screen -r"},
		{"reattach named", "screen -r myproject", "screen -r <val>"},
		{"reattach detach", "screen -d -r myproject", "screen -d -r <val>"},
		{"reattach force", "screen -D -R myproject", "screen -D -R <val>"},
		{"reattach multi", "screen -x shared-session", "screen -x <val>"},

		// Flags with arguments
		{"config file", "screen -c /etc/screenrc", "screen -c <path>"},
		{"shell", "screen -s /bin/zsh", "screen -s <val>"},
		{"terminal type", "screen -T xterm-256color", "screen -T <val>"},
		{"scrollback lines", "screen -h 5000", "screen -h N"},
		{"escape char", "screen -e ^Bb", "screen -e <val>"},
		{"title", "screen -t editor", "screen -t <val>"},
		{"window", "screen -p 2", "screen -p <val>"},

		// Send command (-X is a mode flag; command words are positionals)
		{"send command", "screen -X -S mysession quit", "screen -X -S <val>+"},
		{"send stuff", "screen -S dev -X stuff 'hello world'", "screen -S <val> -X <val>+"},

		// Command to run
		{"run command", "screen vim /etc/hosts", "screen <val> <path>"},
		{"run with session", "screen -S build make -j4", "screen -S <val>+"},

		// Boolean flags
		{"logging", "screen -L -S dev", "screen -L -S <val>"},
		{"version", "screen -v", "screen -v"},
		{"quiet", "screen -q", "screen -q"},

		// Redirect
		{"with redirect", "screen -ls > /tmp/sessions.txt", "screen -ls > <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different session names collide", func(t *testing.T) {
		a := shellshape.Normalize("screen -S project-alpha")
		b := shellshape.Normalize("screen -S project-beta")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different reattach targets collide", func(t *testing.T) {
		a := shellshape.Normalize("screen -r session1")
		b := shellshape.Normalize("screen -r session2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different scrollback sizes collide", func(t *testing.T) {
		a := shellshape.Normalize("screen -h 1000")
		b := shellshape.Normalize("screen -h 9999")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("screen literal-arg")
		subshell := shellshape.Normalize("screen $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
