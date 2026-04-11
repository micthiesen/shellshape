package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestIP(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"addr show", "ip addr show", "ip addr show"},
		{"addr shorthand", "ip a", "ip a"},
		{"link show", "ip link show", "ip link show"},
		{"route show", "ip route show", "ip route show"},
		{"route shorthand", "ip route", "ip route"},
		{"neigh show", "ip neigh show", "ip neigh show"},

		// Global flags
		{"ipv6 route", "ip -6 route", "ip -6 route"},
		{"brief addr", "ip -brief addr show", "ip -brief addr show"},
		{"stats link", "ip -s link show", "ip -s link show"},
		{"json output", "ip -j addr show", "ip -j addr show"},
		{"family inet", "ip -f inet addr show", "ip -f <val> addr show"},
		{"netns flag", "ip -n mynamespace addr show", "ip -n <val> addr show"},

		// addr add/del with CIDR
		{"addr add", "ip addr add 192.168.1.100/24 dev eth0", "ip addr add <addr> dev <val>"},
		{"addr del", "ip addr del 10.0.0.5/32 dev wlan0", "ip addr del <addr> dev <val>"},
		{"address flush dev", "ip addr flush dev eth0", "ip addr flush dev <val>"},

		// link set
		{"link set up", "ip link set dev eth0 up", "ip link set dev <val> up"},
		{"link set down", "ip link set dev wlan0 down", "ip link set dev <val> down"},
		{"link set mtu", "ip link set dev eth0 mtu 9000", "ip link set dev <val> mtu N"},
		{"link set mac", "ip link set dev eth0 address 00:11:22:33:44:55", "ip link set dev <val> address <addr>"},
		{"link set txqueuelen", "ip link set dev eth0 txqueuelen 1000", "ip link set dev <val> txqueuelen N"},

		// route add/del
		{"route add default", "ip route add default via 192.168.1.1", "ip route add default via <addr>"},
		{"route add subnet", "ip route add 192.168.0.0/24 dev eth0", "ip route add <addr> dev <val>"},
		{"route del default", "ip route del default", "ip route del default"},
		{"route add via dev", "ip route add 10.0.0.0/8 via 192.168.1.1 dev eth0", "ip route add <addr> via <addr> dev <val>"},
		{"route add metric", "ip route add default via 10.0.0.1 metric 100", "ip route add default via <addr> metric N"},
		{"route get", "ip route get to 8.8.8.8", "ip route get to <addr>"},
		{"route add table", "ip route add 10.0.0.0/8 via 172.16.0.1 table custom", "ip route add <addr> via <addr> table <val>"},

		// scope and label
		{"addr add scope", "ip addr add 192.168.1.1/24 dev eth0 scope global", "ip addr add <addr> dev <val> scope <val>"},
		{"addr add label", "ip addr add 192.168.1.1/24 dev eth0 label eth0:1", "ip addr add <addr> dev <val> label <val>"},

		// tunnel
		{"tunnel show", "ip tunnel show", "ip tunnel show"},

		// subshell as globalValFlag argument
		{"family subshell", "ip -f $(get-family) addr show", "ip -f $(get-family) addr show"},
		// subshell as valKeyword argument
		{"dev subshell", "ip link set dev $(get-dev) up", "ip link set dev $(get-dev) up"},
		// subshell as numKeyword argument
		{"mtu subshell", "ip link set dev eth0 mtu $(get-mtu)", "ip link set dev <val> mtu $(get-mtu)"},
		// subshell as addrKeyword argument
		{"address subshell", "ip link set dev eth0 address $(get-mac)", "ip link set dev <val> address $(get-mac)"},

		// redirect preserved
		{"with redirect", "ip addr show > /tmp/out.txt", "ip addr show > <path>"},
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
	t.Run("different IPs collide", func(t *testing.T) {
		a := shellshape.Normalize("ip addr add 192.168.1.100/24 dev eth0")
		b := shellshape.Normalize("ip addr add 10.0.0.5/32 dev wlan0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different gateways collide", func(t *testing.T) {
		a := shellshape.Normalize("ip route add default via 192.168.1.1")
		b := shellshape.Normalize("ip route add default via 10.0.0.1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("ip link set dev eth0 up")
		b := shellshape.Normalize("ip link set dev wlan0 up")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ip addr add 192.168.1.1/24 dev eth0")
		subshell := shellshape.Normalize("ip addr add $(dangerous-command) dev eth0")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
