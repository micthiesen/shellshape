package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWpctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"status", "wpctl status", "wpctl status"},
		{"inspect", "wpctl inspect 47", "wpctl inspect <val>"},
		{"get-volume", "wpctl get-volume 47", "wpctl get-volume <val>"},

		// Volume control
		{"set-volume numeric", "wpctl set-volume 47 0.5", "wpctl set-volume <val>+"},
		{"set-volume percent", "wpctl set-volume 47 50%", "wpctl set-volume <val>+"},
		{"set-volume relative", "wpctl set-volume 47 5%+", "wpctl set-volume <val>+"},
		{"set-volume at-default", "wpctl set-volume @DEFAULT_AUDIO_SINK@ 0.8", "wpctl set-volume <val>+"},

		// Mute control
		{"set-mute", "wpctl set-mute 47 1", "wpctl set-mute <val>+"},
		{"set-mute toggle", "wpctl set-mute 47 toggle", "wpctl set-mute <val>+"},

		// Set default
		{"set-default", "wpctl set-default 47", "wpctl set-default <val>"},

		// Set profile
		{"set-profile", "wpctl set-profile 47 0", "wpctl set-profile <val>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test
	t.Run("different node IDs collide", func(t *testing.T) {
		a := shellshape.Normalize("wpctl set-volume 47 0.5")
		b := shellshape.Normalize("wpctl set-volume 83 0.8")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("wpctl set-volume 47 0.5")
		subshell := shellshape.Normalize("wpctl set-volume $(evil-cmd) 0.5")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
