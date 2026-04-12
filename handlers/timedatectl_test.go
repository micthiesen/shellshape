package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTimedatectl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// No-argument subcommands
		{"status", "timedatectl status", "timedatectl status"},
		{"show", "timedatectl show", "timedatectl show"},
		{"list-timezones", "timedatectl list-timezones", "timedatectl list-timezones"},

		// Subcommands with values
		{"set-timezone", "timedatectl set-timezone America/New_York", "timedatectl set-timezone <val>"},
		{"set-time", "timedatectl set-time '2023-01-01 12:00:00'", "timedatectl set-time <val>"},
		{"set-ntp true", "timedatectl set-ntp true", "timedatectl set-ntp <val>"},
		{"set-ntp false", "timedatectl set-ntp false", "timedatectl set-ntp <val>"},

		// Boolean flags
		{"no-pager", "timedatectl --no-pager status", "timedatectl --no-pager status"},
		{"no-ask-password", "timedatectl --no-ask-password set-timezone UTC", "timedatectl --no-ask-password set-timezone <val>"},

		// Flags with values
		{"host flag", "timedatectl -H root@server status", "timedatectl -H <val> status"},
		{"property flag", "timedatectl show -p Timezone", "timedatectl show -p <val>"},
		{"property fused", "timedatectl show --property=Timezone", "timedatectl show --property=<val>"},
		{"value flag", "timedatectl show -p Timezone --value", "timedatectl show -p <val> --value"},
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
	t.Run("different timezones collide", func(t *testing.T) {
		a := shellshape.Normalize("timedatectl set-timezone America/New_York")
		b := shellshape.Normalize("timedatectl set-timezone Europe/London")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different times collide", func(t *testing.T) {
		a := shellshape.Normalize("timedatectl set-time '2023-01-01 12:00:00'")
		b := shellshape.Normalize("timedatectl set-time '2024-06-15 08:30:00'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("timedatectl set-timezone literal-arg")
		subshell := shellshape.Normalize("timedatectl set-timezone $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
