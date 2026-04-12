package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKcmshell(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single module", "kcmshell6 kcm_colors", "kcmshell6 kcm_colors"},
		{"list modules", "kcmshell6 --list", "kcmshell6 --list"},
		{"multiple modules", "kcmshell6 kcm_colors kcm_fonts", "kcmshell6 kcm_colors kcm_fonts"},
		{"args flag", "kcmshell6 --args somevalue kcm_colors", "kcmshell6 --args <val> kcm_colors"},
		{"kcmshell5 alias", "kcmshell5 kcm_colors", "kcmshell5 kcm_colors"},
		{"kcmshell5 list", "kcmshell5 --list", "kcmshell5 --list"},
		{"no args", "kcmshell6", "kcmshell6"},
		{"module with dots", "kcmshell6 kcm_kwin_effects", "kcmshell6 kcm_kwin_effects"},
		{"args with multiple modules", "kcmshell6 --args foo kcm_colors kcm_fonts", "kcmshell6 --args <val> kcm_colors kcm_fonts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different --args values produce same shape
	t.Run("different args values collide", func(t *testing.T) {
		a := shellshape.Normalize("kcmshell6 --args value1 kcm_colors")
		b := shellshape.Normalize("kcmshell6 --args value2 kcm_colors")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kcmshell6 kcm_colors")
		subshell := shellshape.Normalize("kcmshell6 $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
