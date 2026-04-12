package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKioclient(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands with paths
		{"exec path", "kioclient exec /tmp/file.pdf", "kioclient exec <path>"},
		{"cat path", "kioclient cat /home/user/doc.txt", "kioclient cat <path>"},
		{"ls path", "kioclient ls /home/user/Documents", "kioclient ls <path>"},
		{"stat path", "kioclient stat ./file.txt", "kioclient stat <path>"},
		{"mkdir path", "kioclient mkdir /tmp/newdir", "kioclient mkdir <path>"},
		{"remove path", "kioclient remove /tmp/old.txt", "kioclient remove <path>"},

		// Copy and move take two paths (src, dest)
		{"copy two paths", "kioclient copy /tmp/src.txt /tmp/dest.txt", "kioclient copy <path>+"},
		{"move two paths", "kioclient move /home/a/file.txt /home/b/file.txt", "kioclient move <path>+"},

		// URLs
		{"exec URL", "kioclient exec https://example.com", "kioclient exec <https-uri>"},
		{"cat smb URL", "kioclient cat smb://server/share/file.txt", "kioclient cat <smb-uri>"},

		// kioclient6 alias
		{"kioclient6 exec", "kioclient6 exec /tmp/file.pdf", "kioclient6 exec <path>"},
		{"kioclient6 ls", "kioclient6 ls /home/user", "kioclient6 ls <path>"},

		// appmenu (no args typically)
		{"appmenu", "kioclient appmenu", "kioclient appmenu"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("kioclient exec /tmp/file.pdf")
		b := shellshape.Normalize("kioclient exec /home/user/image.png")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kioclient exec /tmp/file.pdf")
		subshell := shellshape.Normalize("kioclient exec $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
