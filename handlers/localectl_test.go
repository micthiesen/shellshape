package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLocalectl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// No-argument subcommands
		{"status", "localectl status", "localectl status"},
		{"list-locales", "localectl list-locales", "localectl list-locales"},
		{"list-keymaps", "localectl list-keymaps", "localectl list-keymaps"},

		// Subcommands with values
		{"set-locale", "localectl set-locale LANG=en_US.UTF-8", "localectl set-locale <val>"},
		{"set-locale multiple", "localectl set-locale LANG=en_US.UTF-8 LC_TIME=de_DE.UTF-8", "localectl set-locale <val>"},
		{"set-keymap", "localectl set-keymap us", "localectl set-keymap <val>"},
		{"set-x11-keymap", "localectl set-x11-keymap us pc105 intl", "localectl set-x11-keymap <val>"},

		// Boolean flags
		{"no-pager", "localectl --no-pager status", "localectl --no-pager status"},
		{"no-ask-password", "localectl --no-ask-password set-keymap de", "localectl --no-ask-password set-keymap <val>"},
		{"no-convert", "localectl set-keymap --no-convert us", "localectl set-keymap --no-convert <val>"},

		// Flags with values
		{"host flag", "localectl -H root@server status", "localectl -H <val> status"},
		{"host long", "localectl --host admin@box set-keymap de", "localectl --host <val> set-keymap <val>"},
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
	t.Run("different locales collide", func(t *testing.T) {
		a := shellshape.Normalize("localectl set-locale LANG=en_US.UTF-8")
		b := shellshape.Normalize("localectl set-locale LANG=de_DE.UTF-8")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different keymaps collide", func(t *testing.T) {
		a := shellshape.Normalize("localectl set-keymap us")
		b := shellshape.Normalize("localectl set-keymap de")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("localectl set-keymap literal-arg")
		subshell := shellshape.Normalize("localectl set-keymap $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
