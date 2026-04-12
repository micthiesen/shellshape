package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLoginctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands without arguments
		{"list-sessions", "loginctl list-sessions", "loginctl list-sessions"},
		{"list-users", "loginctl list-users", "loginctl list-users"},

		// Session subcommands (session IDs → N)
		{"show-session", "loginctl show-session 42", "loginctl show-session N"},
		{"terminate-session", "loginctl terminate-session 7", "loginctl terminate-session N"},

		// User subcommands (usernames → <val>)
		{"show-user", "loginctl show-user michael", "loginctl show-user <val>"},
		{"enable-linger", "loginctl enable-linger michael", "loginctl enable-linger <val>"},
		{"disable-linger", "loginctl disable-linger michael", "loginctl disable-linger <val>"},

		// Boolean flags
		{"no-pager", "loginctl --no-pager list-sessions", "loginctl --no-pager list-sessions"},
		{"all flag", "loginctl show-session 1 --all", "loginctl show-session N --all"},

		// Flags with values
		{"property flag", "loginctl show-session 1 --property ActiveState", "loginctl show-session N --property <val>"},
		{"property fused", "loginctl show-user michael --property=Display", "loginctl show-user <val> --property=<val>"},
		{"host flag", "loginctl -H root@server list-sessions", "loginctl -H <val> list-sessions"},

		// Value flag
		{"value flag", "loginctl show-user michael -p State --value", "loginctl show-user <val> -p <val> --value"},
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
	t.Run("different session IDs collide", func(t *testing.T) {
		a := shellshape.Normalize("loginctl show-session 42")
		b := shellshape.Normalize("loginctl show-session 99")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different usernames collide", func(t *testing.T) {
		a := shellshape.Normalize("loginctl show-user alice")
		b := shellshape.Normalize("loginctl show-user bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("loginctl show-user literal-arg")
		subshell := shellshape.Normalize("loginctl show-user $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
