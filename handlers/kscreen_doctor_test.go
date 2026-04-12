package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKscreenDoctor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list outputs", "kscreen-doctor -o", "kscreen-doctor -o"},
		{"list outputs long", "kscreen-doctor --outputs", "kscreen-doctor --outputs"},
		{"json output", "kscreen-doctor -j", "kscreen-doctor -j"},

		// Output configurations
		{"set mode", "kscreen-doctor output.DP-1.mode.1920x1080@60", "kscreen-doctor <val>"},
		{"enable output", "kscreen-doctor output.DP-1.enable", "kscreen-doctor <val>"},
		{"disable output", "kscreen-doctor output.HDMI-1.disable", "kscreen-doctor <val>"},
		{"set position", "kscreen-doctor output.DP-1.position.0,0", "kscreen-doctor <val>"},
		{"multiple configs", "kscreen-doctor output.DP-1.enable output.DP-2.disable", "kscreen-doctor <val>+"},

		// Combined flags and configs
		{"json with config", "kscreen-doctor -j output.DP-1.mode.2560x1440@144", "kscreen-doctor -j <val>"},

		// Scale and rotation
		{"set scale", "kscreen-doctor output.DP-1.scale.2", "kscreen-doctor <val>"},
		{"set rotation", "kscreen-doctor output.eDP-1.rotation.left", "kscreen-doctor <val>"},
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
	t.Run("different output configs collide", func(t *testing.T) {
		a := shellshape.Normalize("kscreen-doctor output.DP-1.mode.1920x1080@60")
		b := shellshape.Normalize("kscreen-doctor output.HDMI-1.mode.2560x1440@144")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kscreen-doctor literal-arg")
		subshell := shellshape.Normalize("kscreen-doctor $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
