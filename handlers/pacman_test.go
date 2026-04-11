package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPacman(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic sync operations
		{"sync update", "pacman -Syu", "pacman -Syu"},
		{"sync install single", "pacman -S nginx", "pacman -S nginx"},
		{"sync install multiple", "pacman -S nginx vim git", "pacman -S nginx vim git"},
		{"sync with needed", "pacman -S --needed base-devel git", "pacman -S --needed base-devel git"},
		{"sync search", "pacman -Ss docker", "pacman -Ss <query>"},
		{"sync info", "pacman -Si nginx", "pacman -Si nginx"},
		{"sync clean cache", "pacman -Scc", "pacman -Scc"},
		{"sync noconfirm", "pacman -S --noconfirm nginx", "pacman -S --noconfirm nginx"},

		// Query operations
		{"query all", "pacman -Q", "pacman -Q"},
		{"query explicit", "pacman -Qe", "pacman -Qe"},
		{"query orphans", "pacman -Qdt", "pacman -Qdt"},
		{"query search", "pacman -Qs network", "pacman -Qs <query>"},
		{"query info", "pacman -Qi nginx", "pacman -Qi nginx"},
		{"query list files", "pacman -Ql nginx", "pacman -Ql nginx"},
		{"query owns file", "pacman -Qo /usr/bin/vim", "pacman -Qo <path>"},
		{"query quiet orphans", "pacman -Qtdq", "pacman -Qtdq"},

		// Remove operations
		{"remove single", "pacman -R nginx", "pacman -R nginx"},
		{"remove with deps", "pacman -Rs nginx", "pacman -Rs nginx"},
		{"remove cascade", "pacman -Rcs nginx", "pacman -Rcs nginx"},
		{"remove nosave", "pacman -Rn nginx", "pacman -Rn nginx"},

		// Upgrade from file
		{"upgrade local file", "pacman -U /tmp/package-1.0-1-x86_64.pkg.tar.zst", "pacman -U <path>"},
		{"upgrade local url", "pacman -U https://example.com/pkg.tar.zst", "pacman -U <https-uri>"},

		// Database operations
		{"database asdeps", "pacman -D --asdeps nginx", "pacman -D --asdeps nginx"},
		{"database asexplicit", "pacman -D --asexplicit nginx", "pacman -D --asexplicit nginx"},

		// Files operations
		{"files search", "pacman -Fs libxml", "pacman -Fs <query>"},
		{"files list", "pacman -Fl nginx", "pacman -Fl nginx"},
		{"files refresh", "pacman -Fy", "pacman -Fy"},

		// Long-form operation flags
		{"long sync", "pacman --sync nginx", "pacman --sync nginx"},
		{"long query", "pacman --query --explicit", "pacman --query --explicit"},
		{"long remove", "pacman --remove nginx", "pacman --remove nginx"},

		// Global flags with values
		{"config flag", "pacman --config /etc/alt.conf -S nginx", "pacman --config <path> -S nginx"},
		{"dbpath flag", "pacman --dbpath /tmp/db -Q", "pacman --dbpath <path> -Q"},
		{"root flag", "pacman -r /mnt -S base", "pacman -r <path> -S base"},
		{"arch flag", "pacman --arch x86_64 -S nginx", "pacman --arch <val> -S nginx"},
		{"color flag", "pacman --color always -Qi nginx", "pacman --color <val> -Qi nginx"},

		// Aliases
		{"yay install", "yay -S firefox", "yay -S firefox"},
		{"yay search", "yay -Ss browser", "yay -Ss <query>"},
		{"paru install", "paru -S firefox", "paru -S firefox"},
		{"paru search", "paru -Ss browser", "paru -Ss <query>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different search terms collide", func(t *testing.T) {
		a := shellshape.Normalize("pacman -Ss docker")
		b := shellshape.Normalize("pacman -Ss kubernetes")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different local packages collide", func(t *testing.T) {
		a := shellshape.Normalize("pacman -U /tmp/foo-1.0.pkg.tar.zst")
		b := shellshape.Normalize("pacman -U /tmp/bar-2.0.pkg.tar.zst")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pacman -S literal-arg")
		subshell := shellshape.Normalize("pacman -S $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
