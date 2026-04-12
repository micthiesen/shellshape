package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXrandr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"query", "xrandr --query", "xrandr --query"},
		{"verbose query", "xrandr --query --verbose", "xrandr --query --verbose"},
		{"no args", "xrandr", "xrandr"},

		// Setting output mode
		{"set mode", "xrandr --output DP-1 --mode 1920x1080", "xrandr --output <val> --mode <val>"},
		{"set rate", "xrandr --output HDMI-1 --mode 2560x1440 --rate 144", "xrandr --output <val> --mode <val> --rate N"},
		{"turn off output", "xrandr --output DP-2 --off", "xrandr --output <val> --off"},
		{"rotate", "xrandr --output DP-1 --rotate left", "xrandr --output <val> --rotate left"},
		{"rotate normal", "xrandr --output DP-1 --rotate normal", "xrandr --output <val> --rotate normal"},

		// Position
		{"position", "xrandr --output DP-1 --pos 1920x0", "xrandr --output <val> --pos <val>"},
		{"left-of", "xrandr --output HDMI-1 --left-of DP-1", "xrandr --output <val> --left-of <val>"},
		{"right-of", "xrandr --output HDMI-1 --right-of DP-1", "xrandr --output <val> --right-of <val>"},

		// Multi-output
		{"multi output", "xrandr --output DP-1 --mode 2560x1440 --primary --output HDMI-1 --mode 1920x1080 --right-of DP-1",
			"xrandr --output <val> --mode <val> --primary --output <val> --mode <val> --right-of <val>"},

		// Screen and dpi
		{"screen", "xrandr --screen 0 --query", "xrandr --screen N --query"},
		{"dpi", "xrandr --dpi 96", "xrandr --dpi N"},

		// Auto
		{"auto", "xrandr --output DP-1 --auto", "xrandr --output <val> --auto"},

		// Dryrun
		{"dryrun", "xrandr --output DP-1 --mode 1920x1080 --dryrun", "xrandr --output <val> --mode <val> --dryrun"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different output names produce same shape
	t.Run("different outputs collide", func(t *testing.T) {
		a := shellshape.Normalize("xrandr --output DP-1 --mode 1920x1080")
		b := shellshape.Normalize("xrandr --output HDMI-2 --mode 2560x1440")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xrandr --output DP-1")
		subshell := shellshape.Normalize("xrandr --output $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
