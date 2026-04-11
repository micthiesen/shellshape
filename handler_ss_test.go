package shellshape

import "testing"

func TestSs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "ss", "ss"},
		{"boolean flags", "ss -tlp", "ss -tlp"},
		{"listening numeric", "ss -ltn", "ss -ltn"},
		{"all tcp", "ss -t -a", "ss -t -a"},
		{"summary", "ss -s", "ss -s"},

		// State filters
		{"state filter", "ss state LISTEN", "ss state <state>"},
		{"state established", "ss -t state established", "ss -t state <state>"},
		{"exclude state", "ss exclude time-wait", "ss exclude <state>"},

		// Address filters
		{"dst filter", "ss -t dst 192.168.2.1:80", "ss -t dst <addr>"},
		{"src filter", "ss src :443", "ss src <addr>"},
		{"src port", "ss -t src :8080", "ss -t src <addr>"},

		// Quoted filter expressions
		{"sport filter", "ss 'sport = :80'", "ss <filter>"},
		{"dport filter", "ss -t -a 'dport = :22'", "ss -t -a <filter>"},
		{"complex filter", "ss -t '( dport = :ssh or sport = :ssh )'", "ss -t <filter>"},

		// Flags with arguments
		{"family flag", "ss -f inet", "ss -f <val>"},
		{"query flag", "ss -A tcp", "ss -A <val>"},
		{"diag flag", "ss -D /tmp/ss-dump", "ss -D <path>"},
		{"filter file flag", "ss -F /etc/ss-filter", "ss -F <path>"},

		// Combined state and filter
		{"state plus filter", "ss -t -a 'dport = :22' state ESTABLISHED", "ss -t -a <filter> state <state>"},

		// IPv4/IPv6
		{"ipv4", "ss -4 -tlp", "ss -4 -tlp"},
		{"ipv6", "ss -6 -tlp", "ss -6 -tlp"},

		// Redirect
		{"redirect", "ss -tlp > /tmp/out.txt", "ss -tlp > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different addresses collide", func(t *testing.T) {
		a := Normalize("ss -t dst 192.168.1.1:80")
		b := Normalize("ss -t dst 10.0.0.1:443")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different states collide", func(t *testing.T) {
		a := Normalize("ss state LISTEN")
		b := Normalize("ss state ESTABLISHED")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different filters collide", func(t *testing.T) {
		a := Normalize("ss 'sport = :80'")
		b := Normalize("ss 'dport = :443'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ss -t dst 192.168.1.1:80")
		subshell := Normalize("ss -t dst $(echo evil)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
