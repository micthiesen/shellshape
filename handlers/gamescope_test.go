package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGamescope(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"basic with inner command", "gamescope -- game", "gamescope -- <val>"},
		{"fullscreen", "gamescope -f -- game", "gamescope -f -- <val>"},
		{"resolution flags", "gamescope -w 1920 -h 1080 -- game", "gamescope -w N -h N -- <val>"},
		{"output resolution", "gamescope -W 2560 -H 1440 -- game", "gamescope -W N -H N -- <val>"},
		{"refresh rate", "gamescope -r 144 -- game", "gamescope -r N -- <val>"},
		{"hdr and borderless", "gamescope --hdr-enabled -b -- game", "gamescope --hdr-enabled -b -- <val>"},
		{"backend flag", "gamescope --backend wayland -- game", "gamescope --backend wayland -- <val>"},
		{"force grab cursor", "gamescope --force-grab-cursor -- game", "gamescope --force-grab-cursor -- <val>"},
		{"expose steam", "gamescope -e -- game", "gamescope -e -- <val>"},
		{"full gaming setup", "gamescope -w 1920 -h 1080 -W 2560 -H 1440 -r 144 -f --hdr-enabled -- mangohud game --arg", "gamescope -w N -h N -W N -H N -r N -f --hdr-enabled -- <val>+"},
		{"inner command with path", "gamescope -- /usr/bin/game", "gamescope -- <path>"},
		{"no separator bare command", "gamescope game", "gamescope <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different resolutions collide", func(t *testing.T) {
		a := shellshape.Normalize("gamescope -w 1920 -h 1080 -- game")
		b := shellshape.Normalize("gamescope -w 2560 -h 1440 -- game")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different inner commands collide", func(t *testing.T) {
		a := shellshape.Normalize("gamescope -f -- game1")
		b := shellshape.Normalize("gamescope -f -- game2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("gamescope -- game")
		subshell := shellshape.Normalize("gamescope -- $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
