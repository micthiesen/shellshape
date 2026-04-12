package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMakepkg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "makepkg", "makepkg"},
		{"sync deps", "makepkg -s", "makepkg -s"},
		{"sync and install", "makepkg -si", "makepkg -si"},
		{"bundled flags", "makepkg -sicf", "makepkg -sicf"},
		{"long flags", "makepkg --syncdeps --install --clean", "makepkg --syncdeps --install --clean"},
		{"noconfirm", "makepkg --noconfirm -si", "makepkg --noconfirm -si"},
		{"printsrcinfo", "makepkg --printsrcinfo", "makepkg --printsrcinfo"},
		{"verifysource", "makepkg --verifysource", "makepkg --verifysource"},

		// Flags with arguments
		{"pkgbuild short", "makepkg -p custom.PKGBUILD", "makepkg -p <path>"},
		{"config", "makepkg --config /etc/makepkg.conf -si", "makepkg --config <path> -si"},
		{"key", "makepkg --key ABCD1234", "makepkg --key <val>"},

		// Redirect
		{"srcinfo redirect", "makepkg --printsrcinfo > .SRCINFO", "makepkg --printsrcinfo > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different configs collide", func(t *testing.T) {
		a := shellshape.Normalize("makepkg --config /etc/makepkg.conf")
		b := shellshape.Normalize("makepkg --config /home/user/custom.conf")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("makepkg -p myfile.PKGBUILD")
		subshell := shellshape.Normalize("makepkg -p $(echo evil)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
