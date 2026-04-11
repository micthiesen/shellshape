package shellshape

import "testing"

func TestNft(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list ruleset", "nft list ruleset", "nft list ruleset"},
		{"flush ruleset", "nft flush ruleset", "nft flush ruleset"},

		// Table operations (1 name = <word>)
		{"add table", "nft add table inet filter", "nft add table inet <word>"},
		{"delete table", "nft delete table inet myfilter", "nft delete table inet <word>"},
		{"list table", "nft list table inet filter", "nft list table inet <word>"},
		{"add table no family", "nft add table mytable", "nft add table <word>"},

		// Chain operations (2 names = <word>+)
		{"add chain", "nft add chain inet filter input", "nft add chain inet <word>+"},
		{"delete chain", "nft delete chain inet filter forward", "nft delete chain inet <word>+"},
		{"list chain", "nft list chain inet filter input", "nft list chain inet <word>+"},

		// Rule operations - rule body collapses to <val>
		{"add rule", "nft add rule inet filter input tcp dport 22 accept", "nft add rule inet <word>+ <val>"},
		{"add rule complex", "nft add rule inet filter input ip saddr 192.168.0.0/24 accept", "nft add rule inet <word>+ <val>"},
		{"insert rule", "nft insert rule inet filter input position 0 tcp dport 80 accept", "nft insert rule inet <word>+ <val>"},
		{"delete rule handle", "nft delete rule inet filter input handle 3", "nft delete rule inet <word>+ handle N"},

		// Set and element operations
		{"add set", "nft add set inet filter myset", "nft add set inet <word>+"},
		{"add element", "nft add element inet filter myset { 10.0.0.1 }", "nft add element inet <word>+ <val>"},

		// Flags with arguments
		{"file flag short", "nft -f /etc/nftables.conf", "nft -f <path>"},
		{"file flag long", "nft --file /etc/nftables.conf", "nft --file <path>"},
		{"include path", "nft -I /etc/nftables.d", "nft -I <path>"},
		{"define var", "nft -D myvar=foo list ruleset", "nft -D <val> list ruleset"},

		// Boolean flags
		{"numeric and handle", "nft --handle --numeric list chain inet filter input", "nft --handle --numeric list chain inet <word>+"},
		{"json output", "nft -j list ruleset", "nft -j list ruleset"},
		{"check flag", "nft -c -f /etc/nftables.conf", "nft -c -f <path>"},

		// Redirect
		{"list with redirect", "nft list ruleset > /etc/nftables.conf", "nft list ruleset > <path>"},

		// Export/monitor (verbatim)
		{"export json", "nft export json", "nft export json"},
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
	t.Run("different table names collide", func(t *testing.T) {
		a := Normalize("nft add table inet filter")
		b := Normalize("nft add table inet mychain")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different rule bodies collide", func(t *testing.T) {
		a := Normalize("nft add rule inet filter input tcp dport 22 accept")
		b := Normalize("nft add rule inet filter input udp sport 53 drop")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("nft add table inet mytable")
		subshell := Normalize("nft add table inet $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
