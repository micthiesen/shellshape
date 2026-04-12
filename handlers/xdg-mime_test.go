package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXdgMime(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// query filetype
		{"query filetype", "xdg-mime query filetype /tmp/file.pdf", "xdg-mime query filetype <path>"},
		{"query filetype relative", "xdg-mime query filetype ./image.png", "xdg-mime query filetype <path>"},

		// query default
		{"query default mime", "xdg-mime query default image/png", "xdg-mime query default <val>"},
		{"query default text", "xdg-mime query default text/html", "xdg-mime query default <val>"},

		// default (set default app)
		{"default single mime", "xdg-mime default imv.desktop image/png", "xdg-mime default <val>+"},
		{"default multiple mimes", "xdg-mime default imv.desktop image/png image/jpeg", "xdg-mime default <val>+"},

		// install
		{"install xml", "xdg-mime install /usr/share/mime/packages/custom.xml", "xdg-mime install <path>"},
		{"install with novendor", "xdg-mime install --novendor /tmp/myapp.xml", "xdg-mime install --novendor <path>"},
		{"install with mode", "xdg-mime install --mode user ./custom.xml", "xdg-mime install --mode <val> <path>"},

		// uninstall
		{"uninstall xml", "xdg-mime uninstall /usr/share/mime/packages/custom.xml", "xdg-mime uninstall <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different mime types collide", func(t *testing.T) {
		a := shellshape.Normalize("xdg-mime query default image/png")
		b := shellshape.Normalize("xdg-mime query default application/pdf")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("different desktop files collide", func(t *testing.T) {
		a := shellshape.Normalize("xdg-mime default imv.desktop image/png")
		b := shellshape.Normalize("xdg-mime default firefox.desktop text/html")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xdg-mime query default text/html")
		subshell := shellshape.Normalize("xdg-mime query default $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
