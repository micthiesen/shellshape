package shellshape

import "testing"

func TestDig(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple lookup", "dig example.com", "dig <name>"},
		{"with type", "dig example.com MX", "dig <name> MX"},
		{"with server", "dig @8.8.8.8 example.com", "dig @<server> <name>"},
		{"with server and type", "dig @8.8.8.8 example.com AAAA", "dig @<server> <name> AAAA"},
		{"reverse lookup", "dig -x 192.168.1.1", "dig -x <addr>"},

		// Query options (+ prefixed)
		{"short output", "dig +short example.com", "dig +short <name>"},
		{"multiple queryopts", "dig +noall +answer example.com A", "dig +noall +answer <name> A"},
		{"trace", "dig +trace example.com", "dig +trace <name>"},
		{"queryopt with value", "dig +timeout=10 example.com", "dig +timeout=10 <name>"},

		// Flags with arguments
		{"port flag", "dig -p 5353 example.com", "dig -p N <name>"},
		{"query flag", "dig -q example.com", "dig -q <name>"},
		{"type flag", "dig -t MX example.com", "dig -t MX <name>"},
		{"class flag", "dig -c IN example.com", "dig -c IN <name>"},
		{"batch file", "dig -f queries.txt", "dig -f <path>"},
		{"key file", "dig -k tsig.key example.com", "dig -k <path> <name>"},
		{"bind address", "dig -b 10.0.0.1 example.com", "dig -b <addr> <name>"},
		{"tsig key", "dig -y hmac-sha256:mykey:secret== example.com", "dig -y <key> <name>"},

		// Boolean flags
		{"ipv4 only", "dig -4 example.com", "dig -4 <name>"},
		{"ipv6 only", "dig -6 example.com", "dig -6 <name>"},

		// Combined
		{"server port type queryopt", "dig @ns1.example.com -p 53 example.com ANY +noall +answer", "dig @<server> -p N <name> ANY +noall +answer"},

		// Record types preserved
		{"SOA type", "dig example.com SOA", "dig <name> SOA"},
		{"TXT type", "dig example.com TXT", "dig <name> TXT"},
		{"AXFR type", "dig @ns1.example.com example.com AXFR", "dig @<server> <name> AXFR"},
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
	t.Run("different domains collide", func(t *testing.T) {
		a := Normalize("dig example.com MX")
		b := Normalize("dig google.com MX")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different servers collide", func(t *testing.T) {
		a := Normalize("dig @8.8.8.8 example.com")
		b := Normalize("dig @1.1.1.1 google.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("dig example.com")
		subshell := Normalize("dig $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
