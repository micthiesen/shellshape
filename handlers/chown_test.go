package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestChown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"user only", "chown root file.txt", "chown root <path>"},
		{"user and group", "chown root:wheel /etc/config", "chown root:wheel <path>"},
		{"group only", "chown :staff report.pdf", "chown :staff <path>"},
		{"numeric ids", "chown 1000:1000 /app", "chown 1000:1000 <path>"},

		// Flags
		{"recursive", "chown -R www-data:www-data /var/www", "chown -R www-data:www-data <path>"},
		{"verbose recursive", "chown -Rv alice src/", "chown -Rv alice <path>"},
		{"force", "chown -f root /etc/passwd", "chown -f root <path>"},
		{"symlink", "chown -h nobody link.txt", "chown -h nobody <path>"},

		// Multiple files
		{"multiple files", "chown alice a.txt b.txt c.txt", "chown alice <path>+"},
		{"multiple with flags", "chown -R deploy:deploy src/ lib/ bin/", "chown -R deploy:deploy <path>+"},

		// Long flags (GNU)
		{"long recursive", "chown --recursive admin /opt/app", "chown --recursive admin <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different file paths collapse to same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("chown root:wheel /home/alice/file.txt")
		b := shellshape.Normalize("chown root:wheel /opt/deploy/config.yml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different owners stay distinct
	t.Run("different owners distinct", func(t *testing.T) {
		a := shellshape.Normalize("chown alice file.txt")
		b := shellshape.Normalize("chown bob file.txt")
		if a == b {
			t.Errorf("different owners should produce different shapes, both got %q", a)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("chown root file.txt")
		subshell := shellshape.Normalize("chown $(dangerous-command) file.txt")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
