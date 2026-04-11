package shellshape

import "testing"

func TestIptables(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list all", "iptables -L", "iptables -L"},
		{"list chain", "iptables -L INPUT", "iptables -L INPUT"},
		{"list verbose", "iptables -vnL", "iptables -vnL"},
		{"list with line numbers", "iptables -L INPUT --line-numbers", "iptables -L INPUT --line-numbers"},

		// Append rules
		{"append accept", "iptables -A INPUT -p tcp --dport 80 -j ACCEPT", "iptables -A INPUT -p tcp --dport <val> -j ACCEPT"},
		{"append drop", "iptables -A INPUT -p tcp --dport 22 -j DROP", "iptables -A INPUT -p tcp --dport <val> -j DROP"},
		{"append with source", "iptables -A INPUT -s 192.168.1.0/24 -p tcp --dport 443 -j ACCEPT", "iptables -A INPUT -s <addr> -p tcp --dport <val> -j ACCEPT"},
		{"append with interface", "iptables -A INPUT -i eth0 -p tcp --dport 902 -j REJECT --reject-with icmp-port-unreachable", "iptables -A INPUT -i <val> -p tcp --dport <val> -j REJECT --reject-with icmp-port-unreachable"},

		// Delete and insert
		{"delete by number", "iptables -D INPUT 2", "iptables -D INPUT N"},
		{"insert at position", "iptables -I INPUT 1 -p tcp --dport 21 -s 10.0.0.1 -j ACCEPT", "iptables -I INPUT N -p tcp --dport <val> -s <addr> -j ACCEPT"},

		// NAT table
		{"nat masquerade", "iptables -t nat -A POSTROUTING -s 192.168.0.0/24 -j MASQUERADE", "iptables -t nat -A POSTROUTING -s <addr> -j MASQUERADE"},
		{"nat dnat", "iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 10.0.0.2:80", "iptables -t nat -A PREROUTING -p tcp --dport <val> -j DNAT --to-destination <val>"},

		// Match modules
		{"state match", "iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT", "iptables -A INPUT -m state --state <val> -j ACCEPT"},
		{"comment match", "iptables -A INPUT -p tcp --dport 80 -j ACCEPT -m comment --comment allow-http", "iptables -A INPUT -p tcp --dport <val> -j ACCEPT -m comment --comment <str>"},
		{"limit match", "iptables -A INPUT -p icmp --icmp-type echo-request -m limit --limit 1/s -j ACCEPT", "iptables -A INPUT -p icmp --icmp-type <val> -m limit --limit <val> -j ACCEPT"},

		// Flush and policy
		{"flush", "iptables -F", "iptables -F"},
		{"flush chain", "iptables -F INPUT", "iptables -F INPUT"},
		{"policy", "iptables -P INPUT DROP", "iptables -P INPUT DROP"},
		{"new chain", "iptables -N MYCHAIN", "iptables -N MYCHAIN"},
		{"delete chain", "iptables -X MYCHAIN", "iptables -X MYCHAIN"},

		// ip6tables alias
		{"ip6tables basic", "ip6tables -A INPUT -p tcp --dport 443 -j DROP", "ip6tables -A INPUT -p tcp --dport <val> -j DROP"},
		{"ip6tables list", "ip6tables -L -n", "ip6tables -L -n"},

		// Sport
		{"source port", "iptables -A OUTPUT -p tcp --sport 1024:65535 -j ACCEPT", "iptables -A OUTPUT -p tcp --sport <val> -j ACCEPT"},

		// Destination and source together
		{"src and dst", "iptables -A FORWARD -s 10.0.0.0/8 -d 172.16.0.0/12 -j ACCEPT", "iptables -A FORWARD -s <addr> -d <addr> -j ACCEPT"},

		// Subshell as standalone token in args
		{"subshell standalone", "iptables $(get-flags) -j ACCEPT", "iptables $(get-flags) -j ACCEPT"},
		// Subshell as verbatim flag argument
		{"subshell verbatim arg", "iptables -A INPUT -p $(get-proto) -j ACCEPT", "iptables -A INPUT -p $(get-proto) -j ACCEPT"},
		// Subshell as val flag argument
		{"subshell val arg", "iptables -A INPUT --dport $(get-port) -j ACCEPT", "iptables -A INPUT --dport $(get-port) -j ACCEPT"},
		// Subshell as str flag argument
		{"subshell str arg", "iptables -A INPUT -m comment --comment $(gen-comment) -j ACCEPT", "iptables -A INPUT -m comment --comment $(gen-comment) -j ACCEPT"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different IPs collide", func(t *testing.T) {
		a := Normalize("iptables -A INPUT -s 10.0.0.1 -j DROP")
		b := Normalize("iptables -A INPUT -s 172.16.5.99 -j DROP")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := Normalize("iptables -A INPUT -p tcp --dport 80 -j ACCEPT")
		b := Normalize("iptables -A INPUT -p tcp --dport 443 -j ACCEPT")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("iptables -A INPUT -s 10.0.0.1 -j DROP")
		subshell := Normalize("iptables -A INPUT -s $(get-ip) -j DROP")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
