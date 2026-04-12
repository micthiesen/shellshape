package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestResolvectl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"status", "resolvectl status", "resolvectl status"},
		{"statistics", "resolvectl statistics", "resolvectl statistics"},
		{"flush-caches", "resolvectl flush-caches", "resolvectl flush-caches"},

		// Query with domains/IPs
		{"query domain", "resolvectl query example.com", "resolvectl query <val>"},
		{"query ip", "resolvectl query 8.8.8.8", "resolvectl query <val>"},
		{"query multiple", "resolvectl query example.com google.com", "resolvectl query <val>+"},

		// Status with interface
		{"status iface", "resolvectl status eth0", "resolvectl status <val>"},

		// DNS configuration
		{"dns set", "resolvectl dns eth0 8.8.8.8 8.8.4.4", "resolvectl dns <val>+"},
		{"domain set", "resolvectl domain eth0 example.com", "resolvectl domain <val>+"},

		// LLMNR/mDNS (interface + yes/no collapses to <val>+)
		{"llmnr", "resolvectl llmnr eth0 yes", "resolvectl llmnr <val>+"},
		{"mdns", "resolvectl mdns eth0 yes", "resolvectl mdns <val>+"},

		// Flags before subcommand (fused flags stay as normalizer emits them)
		{"query with type fused", "resolvectl --type=MX query example.com", "resolvectl --type=MX query <val>"},
		{"query with legend fused", "resolvectl --legend=no query example.com", "resolvectl --legend=no query <val>"},
		{"query short type", "resolvectl -t MX query example.com", "resolvectl -t <val> query <val>"},
		{"interface flag", "resolvectl -i 2 query example.com", "resolvectl -i N query <val>"},
		{"protocol flag", "resolvectl -p dns query example.com", "resolvectl -p <val> query <val>"},

		// Flags after subcommand
		{"query type after", "resolvectl query --type=MX example.com", "resolvectl query --type=<val> <val>"},

		// Service resolution
		{"service", "resolvectl service _xmpp-server._tcp gmail.com", "resolvectl service <val>+"},

		// openpgp/tlsa
		{"openpgp", "resolvectl openpgp user@example.com", "resolvectl openpgp <val>"},
		{"tlsa", "resolvectl tlsa tcp example.com:443", "resolvectl tlsa <val>+"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different domains -> same shape
	t.Run("different domains collide", func(t *testing.T) {
		a := shellshape.Normalize("resolvectl query example.com")
		b := shellshape.Normalize("resolvectl query google.com")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("resolvectl query example.com")
		subshell := shellshape.Normalize("resolvectl query $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
