package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKreadconfig(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"read key", "kreadconfig6 --group General --key loginMode", "kreadconfig6 --group <val> --key <val>"},
		{"read with file", "kreadconfig6 --file ksmserverrc --group General --key loginMode", "kreadconfig6 --file <val> --group <val> --key <val>"},
		{"nested groups", "kreadconfig6 --file kwinrc --group Plugins --group kwin4_effect_blurEnabled --key enabled", "kreadconfig6 --file <val> --group <val> --group <val> --key <val>"},

		// kwriteconfig6
		{"write value flag", "kwriteconfig6 --file kwinrc --group General --key foo --value bar", "kwriteconfig6 --file <val> --group <val> --key <val> --value <val>"},
		{"write with type", "kwriteconfig6 --file kwinrc --group Windows --key BorderlessMaximizedWindows --type bool --value true", "kwriteconfig6 --file <val> --group <val> --key <val> --type bool --value <val>"},
		{"write delete", "kwriteconfig6 --group General --key oldKey --delete", "kwriteconfig6 --group <val> --key <val> --delete"},

		// kreadconfig5 and kwriteconfig5
		{"kreadconfig5", "kreadconfig5 --file startkderc --group General --key systemdBoot", "kreadconfig5 --file <val> --group <val> --key <val>"},
		{"kwriteconfig5", "kwriteconfig5 --file startkderc --group General --key systemdBoot --value true", "kwriteconfig5 --file <val> --group <val> --key <val> --value <val>"},

		// Positional value (some invocations pass value as trailing positional)
		{"trailing positional", "kwriteconfig6 --group General --key foo bar", "kwriteconfig6 --group <val> --key <val>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files and keys collide", func(t *testing.T) {
		a := shellshape.Normalize("kreadconfig6 --file kwinrc --group General --key foo")
		b := shellshape.Normalize("kreadconfig6 --file ksmserverrc --group Plugins --key bar")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kreadconfig6 --file kwinrc --group General --key foo")
		subshell := shellshape.Normalize("kreadconfig6 --file $(dangerous) --group General --key foo")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
