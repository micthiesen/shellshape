package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestXdgOpen(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare command", "xdg-open", "xdg-open"},
		{"open URL", "xdg-open https://example.com", "xdg-open <https-uri>"},
		{"open http URL", "xdg-open http://example.com/page", "xdg-open <http-uri>"},
		{"open current dir", "xdg-open .", "xdg-open ."},
		{"open absolute path", "xdg-open /tmp/file.pdf", "xdg-open <path>"},
		{"open relative path", "xdg-open ./image.png", "xdg-open <path>"},
		{"open file by extension", "xdg-open document.pdf", "xdg-open <path>"},
		{"open parent-relative path", "xdg-open ../file.txt", "xdg-open <path>"},
		{"open path with dirs", "xdg-open src/main.go", "xdg-open <path>"},
		{"help flag", "xdg-open --help", "xdg-open --help"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different URLs collide", func(t *testing.T) {
		a := shellshape.Normalize("xdg-open https://google.com")
		b := shellshape.Normalize("xdg-open https://github.com/foo/bar")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("xdg-open /home/user/doc.pdf")
		b := shellshape.Normalize("xdg-open /tmp/image.png")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("xdg-open https://example.com")
		subshell := shellshape.Normalize("xdg-open $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
