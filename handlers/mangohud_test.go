package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMangohud(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"wrap command", "mangohud game", "mangohud game"},
		{"wrap command with args", "mangohud game --fullscreen", "mangohud game --fullscreen"},
		{"wrap path command", "mangohud /usr/bin/game", "mangohud /usr/bin/game"},
		{"dlsym flag", "mangohud --dlsym game", "mangohud --dlsym game"},
		{"dlsym with args", "mangohud --dlsym game --windowed 1920", "mangohud --dlsym game --windowed N"},
		{"wrap with path args", "mangohud game /path/to/config.ini", "mangohud game <path>"},
		{"wrap with numeric args", "mangohud game 42", "mangohud game N"},
		{"bare mangohud", "mangohud", "mangohud"},
		{"wrap with url arg", "mangohud game https://example.com", "mangohud game <https-uri>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different inner commands collide paths", func(t *testing.T) {
		a := shellshape.Normalize("mangohud game /path/to/save1")
		b := shellshape.Normalize("mangohud game /path/to/save2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mangohud game")
		subshell := shellshape.Normalize("mangohud $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
