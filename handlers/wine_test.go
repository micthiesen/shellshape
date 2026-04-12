package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"exe file", "wine program.exe", "wine <path>"},
		{"unix path exe", "wine '/home/user/.wine/drive_c/game.exe'", "wine <path>"},
		{"windows path", `wine 'C:\windows\system32\notepad.exe'`, "wine <path>"},
		{"exe with flag args", "wine game.exe --fullscreen", "wine <path> <val>"},
		{"exe with installer switch", "wine setup.exe /S", "wine <path> <val>"},
		{"version flag", "wine --version", "wine --version"},
		{"wine64 alias", "wine64 program.exe", "wine64 <path>"},
		{"exe with numeric arg", "wine game.exe -width 1920", "wine <path> <val> N"},
		{"explorer", "wine explorer.exe", "wine <path>"},
		{"exe with string arg", "wine game.exe hello", "wine <path> <val>"},
		{"exe with multiple args", "wine game.exe hello world", "wine <path> <val>+"},
		{"exe with url arg", "wine browser.exe https://example.com", "wine <path> <https-uri>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different exe paths collide", func(t *testing.T) {
		a := shellshape.Normalize("wine game1.exe")
		b := shellshape.Normalize("wine game2.exe")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("wine program.exe")
		subshell := shellshape.Normalize("wine $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
