package shellshape

import "testing"

func TestHost(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple lookup", "host example.com", "host <name>"},
		{"reverse lookup", "host 192.168.1.1", "host <name>"},
		{"with server", "host example.com 8.8.8.8", "host <name> <server>"},

		// Flags with arguments
		{"type flag", "host -t MX example.com", "host -t MX <name>"},
		{"type AAAA with server", "host -t AAAA example.com 1.1.1.1", "host -t AAAA <name> <server>"},
		{"class flag", "host -c CH example.com", "host -c CH <name>"},
		{"retries flag", "host -R 3 example.com", "host -R N <name>"},
		{"timeout flag", "host -W 10 example.com", "host -W N <name>"},
		{"ndots flag", "host -N 2 example.com", "host -N N <name>"},
		{"memory debug flag", "host -m record example.com", "host -m record <name>"},

		// Boolean flags
		{"all flag", "host -a example.com", "host -a <name>"},
		{"verbose flag", "host -v example.com", "host -v <name>"},
		{"ipv4 only", "host -4 example.com", "host -4 <name>"},
		{"ipv6 only", "host -6 example.com", "host -6 <name>"},
		{"tcp flag", "host -T example.com", "host -T <name>"},
		{"list zone", "host -l example.com ns1.example.com", "host -l <name> <server>"},

		// Combined flags
		{"verbose with type", "host -v -t NS example.com", "host -v -t NS <name>"},
		{"multiple numeric flags", "host -R 3 -W 5 example.com", "host -R N -W N <name>"},
		{"all flags combo", "host -a -4 example.com 8.8.4.4", "host -a -4 <name> <server>"},
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
	t.Run("different hostnames collide", func(t *testing.T) {
		a := Normalize("host example.com")
		b := Normalize("host google.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different servers collide", func(t *testing.T) {
		a := Normalize("host example.com 8.8.8.8")
		b := Normalize("host example.com 1.1.1.1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric values collide", func(t *testing.T) {
		a := Normalize("host -R 3 -W 5 example.com")
		b := Normalize("host -R 10 -W 30 google.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("host literal-arg")
		subshell := Normalize("host $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
