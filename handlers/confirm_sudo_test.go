package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestConfirmSudo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - delegates to pacman handler
		{"pacman install", "~/.confirm-sudo.sh pacman -S cachyos-znver4/go", ".confirm-sudo.sh pacman -S cachyos-znver4/go"},
		{"pacman sync", "~/.confirm-sudo.sh pacman -Sy", ".confirm-sudo.sh pacman -Sy"},
		{"pacman noconfirm", "~/.confirm-sudo.sh pacman -S --noconfirm cachyos-znver4/go", ".confirm-sudo.sh pacman -S --noconfirm cachyos-znver4/go"},
		{"pacman remove", "~/.confirm-sudo.sh pacman -R vim", ".confirm-sudo.sh pacman -R vim"},
		{"pacman search", "~/.confirm-sudo.sh pacman -Ss kernel", ".confirm-sudo.sh pacman -Ss <query>"},

		// Delegates to systemctl handler
		{"systemctl start", "~/.confirm-sudo.sh systemctl start nginx", ".confirm-sudo.sh systemctl start <unit>"},
		{"systemctl enable", "~/.confirm-sudo.sh systemctl enable sshd", ".confirm-sudo.sh systemctl enable <unit>"},

		// Simple commands without specific handlers
		{"simple command", "~/.confirm-sudo.sh ls /home/user", ".confirm-sudo.sh ls <path>"},
		{"command with flags", "~/.confirm-sudo.sh rm -rf /tmp/cache", ".confirm-sudo.sh rm -rf <path>"},

		// Absolute path to wrapper
		{"/home path", "/home/michael/.confirm-sudo.sh pacman -Syu", ".confirm-sudo.sh pacman -Syu"},

		// Without tilde prefix (direct name)
		{"direct name", "confirm-sudo.sh pacman -S vim", "confirm-sudo.sh pacman -S vim"},

		// No arguments
		{"no args", "~/.confirm-sudo.sh", ".confirm-sudo.sh"},
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
	t.Run("different packages collide for search", func(t *testing.T) {
		a := shellshape.Normalize("~/.confirm-sudo.sh pacman -Ss linux")
		b := shellshape.Normalize("~/.confirm-sudo.sh pacman -Ss kernel")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("~/.confirm-sudo.sh cat /etc/shadow")
		b := shellshape.Normalize("~/.confirm-sudo.sh cat /etc/passwd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("~/.confirm-sudo.sh literal-arg")
		subshell := shellshape.Normalize("~/.confirm-sudo.sh $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
