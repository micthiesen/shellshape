package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTmux(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands (no args)
		{"list-sessions", "tmux list-sessions", "tmux list-sessions"},
		{"ls alias", "tmux ls", "tmux ls"},

		// new-session with flags
		{"new-session with name", "tmux new-session -s myproject", "tmux new-session -s <val>"},
		{"new alias", "tmux new -s dev", "tmux new -s <val>"},
		{"new-session with dir", "tmux new-session -s dev -c ~/projects/app", "tmux new-session -s <val> -c <path>"},
		{"new-session detached", "tmux new-session -d -s bg-session", "tmux new-session -d -s <val>"},
		{"new-session with window name", "tmux new-session -s main -n editor", "tmux new-session -s <val> -n <val>"},

		// attach
		{"attach with target", "tmux attach -t mysession", "tmux attach -t <val>"},
		{"attach alias", "tmux a -t 0", "tmux a -t <val>"},
		{"attach detach others", "tmux attach -d -t work", "tmux attach -d -t <val>"},

		// kill-session
		{"kill-session", "tmux kill-session -t old-session", "tmux kill-session -t <val>"},

		// send-keys: positionals collapse to <str>
		{"send-keys simple", "tmux send-keys -t dev 'ls -la' Enter", "tmux send-keys -t <val> <str>+"},
		{"send-keys multiple", "tmux send-keys -t dev:0 'echo hello' C-m", "tmux send-keys -t <val> <str>+"},
		{"send-keys single key", "tmux send-keys -t 0 Enter", "tmux send-keys -t <val> <str>"},

		// split-window / select-pane
		{"split-window with dir", "tmux split-window -c ~/code", "tmux split-window -c <path>"},
		{"split-window horizontal", "tmux split-window -h -t dev", "tmux split-window -h -t <val>"},
		{"select-pane", "tmux select-pane -t 2", "tmux select-pane -t <val>"},

		// Format flag
		{"list with format", "tmux list-sessions -F '#{session_name}'", "tmux list-sessions -F <val>"},

		// Redirect
		{"with redirect", "tmux ls > /tmp/sessions.txt", "tmux ls > <path>"},
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
		a := shellshape.Normalize("tmux new-session -s project-alpha")
		b := shellshape.Normalize("tmux new-session -s project-beta")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different targets collide", func(t *testing.T) {
		a := shellshape.Normalize("tmux attach -t session1")
		b := shellshape.Normalize("tmux attach -t session2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different send-keys values collide", func(t *testing.T) {
		a := shellshape.Normalize("tmux send-keys -t dev 'make build' Enter")
		b := shellshape.Normalize("tmux send-keys -t dev 'npm test' C-m")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("tmux send-keys -t dev literal-arg")
		subshell := shellshape.Normalize("tmux send-keys -t dev $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
