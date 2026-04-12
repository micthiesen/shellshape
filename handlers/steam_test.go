package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSteam(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare steam", "steam", "steam"},
		{"steam url", "steam steam://run/730", "steam <steam-uri>"},
		{"steam rungameid", "steam steam://rungameid/440", "steam <steam-uri>"},

		// Boolean flags
		{"silent", "steam -silent", "steam -silent"},
		{"shutdown", "steam -shutdown", "steam -shutdown"},
		{"bigpicture", "steam -bigpicture", "steam -bigpicture"},
		{"console", "steam -console", "steam -console"},

		// App launch with numeric ID
		{"applaunch", "steam -applaunch 730", "steam -applaunch N"},
		{"applaunch with args", "steam -applaunch 440 -novid -windowed", "steam -applaunch N -novid -windowed"},

		// Combined flags
		{"silent bigpicture", "steam -silent -bigpicture", "steam -silent -bigpicture"},
		{"applaunch silent", "steam -silent -applaunch 570", "steam -silent -applaunch N"},

		// With redirect
		{"redirect stderr", "steam -silent 2>/dev/null", "steam -silent 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different app ids collide", func(t *testing.T) {
		a := shellshape.Normalize("steam -applaunch 730")
		b := shellshape.Normalize("steam -applaunch 440")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different steam urls collide", func(t *testing.T) {
		a := shellshape.Normalize("steam steam://run/730")
		b := shellshape.Normalize("steam steam://rungameid/440")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("steam literal-arg")
		subshell := shellshape.Normalize("steam $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
