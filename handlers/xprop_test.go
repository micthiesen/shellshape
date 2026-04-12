package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXprop(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args click mode", "xprop", "xprop"},
		{"root window", "xprop -root", "xprop -root"},
		{"single property", "xprop WM_CLASS", "xprop WM_CLASS"},
		{"multiple properties", "xprop WM_CLASS _NET_WM_NAME", "xprop WM_CLASS _NET_WM_NAME"},
		{"window id", "xprop -id 0x1400007", "xprop -id <val>"},
		{"window id with property", "xprop -id 0x1400007 WM_CLASS", "xprop -id <val> WM_CLASS"},
		{"name flag", "xprop -name Firefox", "xprop -name <val>"},
		{"spy mode", "xprop -root -spy _NET_ACTIVE_WINDOW", "xprop -root -spy _NET_ACTIVE_WINDOW"},
		{"display flag", "xprop -display :0 -root", "xprop -display <val> -root"},
		{"root with multiple props", "xprop -root _NET_WM_NAME _NET_ACTIVE_WINDOW", "xprop -root _NET_WM_NAME _NET_ACTIVE_WINDOW"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different window IDs produce same shape
	t.Run("different window ids collide", func(t *testing.T) {
		a := shellshape.Normalize("xprop -id 0x1400007 WM_CLASS")
		b := shellshape.Normalize("xprop -id 0x2800003 WM_CLASS")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xprop -id 0x1400007")
		subshell := shellshape.Normalize("xprop -id $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
