package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestHostnamectl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// No-argument subcommands
		{"status", "hostnamectl status", "hostnamectl status"},
		{"bare", "hostnamectl", "hostnamectl"},

		// Subcommands that set values
		{"hostname", "hostnamectl hostname myhost", "hostnamectl hostname <val>"},
		{"icon-name", "hostnamectl icon-name computer-desktop", "hostnamectl icon-name <val>"},
		{"chassis", "hostnamectl chassis desktop", "hostnamectl chassis <val>"},
		{"deployment", "hostnamectl deployment production", "hostnamectl deployment <val>"},
		{"location", "hostnamectl location 'Rack 5, Slot 3'", "hostnamectl location <val>"},

		// Read without setting (no positional)
		{"hostname read", "hostnamectl hostname", "hostnamectl hostname"},
		{"chassis read", "hostnamectl chassis", "hostnamectl chassis"},

		// Boolean flags
		{"static flag", "hostnamectl hostname --static myhost", "hostnamectl hostname --static <val>"},
		{"transient flag", "hostnamectl hostname --transient myhost", "hostnamectl hostname --transient <val>"},
		{"pretty flag", "hostnamectl hostname --pretty 'My Pretty Host'", "hostnamectl hostname --pretty <val>"},
		{"no-ask-password", "hostnamectl --no-ask-password hostname myhost", "hostnamectl --no-ask-password hostname <val>"},

		// Flags with values
		{"host flag", "hostnamectl -H root@server status", "hostnamectl -H <val> status"},
		{"host long", "hostnamectl --host admin@box hostname", "hostnamectl --host <val> hostname"},
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
	t.Run("different hostnames collide", func(t *testing.T) {
		a := shellshape.Normalize("hostnamectl hostname myhost")
		b := shellshape.Normalize("hostnamectl hostname otherhost")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different chassis values collide", func(t *testing.T) {
		a := shellshape.Normalize("hostnamectl chassis desktop")
		b := shellshape.Normalize("hostnamectl chassis laptop")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("hostnamectl hostname literal-arg")
		subshell := shellshape.Normalize("hostnamectl hostname $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
