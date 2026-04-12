package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNetworkctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands (no args)
		{"list", "networkctl list", "networkctl list"},
		{"lldp", "networkctl lldp", "networkctl lldp"},
		{"label", "networkctl label", "networkctl label"},

		// Status with interface
		{"status no arg", "networkctl status", "networkctl status"},
		{"status iface", "networkctl status eth0", "networkctl status <val>"},
		{"status multiple", "networkctl status eth0 wlan0", "networkctl status <val>+"},

		// Up/down with interface
		{"up iface", "networkctl up eth0", "networkctl up <val>"},
		{"down iface", "networkctl down wlan0", "networkctl down <val>"},

		// Flags
		{"no-pager", "networkctl --no-pager status", "networkctl --no-pager status"},
		{"no-legend", "networkctl --no-legend list", "networkctl --no-legend list"},
		{"all flag", "networkctl -a list", "networkctl -a list"},
		{"json output", "networkctl --json=short status eth0", "networkctl --json=short status <val>"},

		// List with pattern
		{"list pattern", "networkctl list eth*", "networkctl list <val>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different interfaces -> same shape
	t.Run("different interfaces collide", func(t *testing.T) {
		a := shellshape.Normalize("networkctl status eth0")
		b := shellshape.Normalize("networkctl status wlan0")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("networkctl status eth0")
		subshell := shellshape.Normalize("networkctl status $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
