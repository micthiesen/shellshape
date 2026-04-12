package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNmcli(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"general status", "nmcli general", "nmcli general"},
		{"general status explicit", "nmcli general status", "nmcli general status"},
		{"networking on", "nmcli networking on", "nmcli networking on"},
		{"networking off", "nmcli networking off", "nmcli networking off"},
		{"radio wifi", "nmcli radio wifi", "nmcli radio wifi"},
		{"radio wifi on", "nmcli radio wifi on", "nmcli radio wifi on"},
		{"monitor", "nmcli monitor", "nmcli monitor"},

		// Connection subcommands
		{"con show", "nmcli connection show", "nmcli connection show"},
		{"con show active", "nmcli connection show --active", "nmcli connection show --active"},
		{"con show by name", "nmcli connection show MyWifi", "nmcli connection show <val>"},
		{"con up", "nmcli connection up MyVPN", "nmcli connection up <val>"},
		{"con down", "nmcli connection down eth0-conn", "nmcli connection down <val>"},
		{"con delete", "nmcli connection delete my-connection", "nmcli connection delete <val>"},
		{"con add ethernet", "nmcli connection add type ethernet con-name eth0-static ifname eth0", "nmcli connection add type ethernet con-name <val> ifname <val>"},
		{"con modify props", "nmcli connection modify my-conn ipv4.addresses 192.168.1.100/24", "nmcli connection modify <val>+"},
		{"con modify dns", "nmcli connection modify my-conn ipv4.dns 8.8.8.8", "nmcli connection modify <val>+"},

		// con shorthand
		{"con shorthand show", "nmcli con show", "nmcli con show"},
		{"con shorthand up", "nmcli con up MyVPN", "nmcli con up <val>"},

		// Device subcommands
		{"dev status", "nmcli device status", "nmcli device status"},
		{"dev show", "nmcli device show eth0", "nmcli device show <val>"},
		{"dev wifi list", "nmcli device wifi list", "nmcli device wifi list"},
		{"dev wifi rescan", "nmcli device wifi rescan", "nmcli device wifi rescan"},
		{"dev wifi connect", "nmcli device wifi connect HomeSSID password secret123 ifname wlan0", "nmcli device wifi connect <val> password <val> ifname <val>"},
		{"dev shorthand", "nmcli dev status", "nmcli dev status"},
		{"dev shorthand wifi", "nmcli dev wifi list", "nmcli dev wifi list"},

		// Global flags
		{"terse output", "nmcli -t connection show", "nmcli -t connection show"},
		{"pretty output", "nmcli -p device status", "nmcli -p device status"},
		{"fields flag", "nmcli -f NAME,TYPE connection show", "nmcli -f <val> connection show"},

		// Import/export
		{"con import", "nmcli connection import type openvpn file /etc/openvpn/client.ovpn", "nmcli connection import type openvpn file <path>"},

		// Agent
		{"agent secret", "nmcli agent secret", "nmcli agent secret"},
		{"agent polkit", "nmcli agent polkit", "nmcli agent polkit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different connection names -> same shape
	t.Run("different connection names collide", func(t *testing.T) {
		a := shellshape.Normalize("nmcli connection up HomeWifi")
		b := shellshape.Normalize("nmcli connection up OfficeVPN")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different SSIDs collide", func(t *testing.T) {
		a := shellshape.Normalize("nmcli device wifi connect NetworkA password pass1 ifname wlan0")
		b := shellshape.Normalize("nmcli device wifi connect NetworkB password pass2 ifname wlan1")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nmcli connection up MyVPN")
		subshell := shellshape.Normalize("nmcli connection up $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
