package shellshape

import "testing"

func TestNslookup(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple lookup", "nslookup example.com", "nslookup <name>"},
		{"reverse lookup", "nslookup 192.168.1.1", "nslookup <name>"},
		{"with server", "nslookup example.com 8.8.8.8", "nslookup <name> <server>"},

		// Flags with type values (preserved verbatim)
		{"type flag", "nslookup -type=MX example.com", "nslookup -type=MX <name>"},
		{"query flag", "nslookup -query=ns example.com", "nslookup -query=ns <name>"},
		{"querytype flag", "nslookup -querytype=AAAA example.com", "nslookup -querytype=AAAA <name>"},
		{"type PTR reverse", "nslookup -type=PTR 54.240.162.118", "nslookup -type=PTR <name>"},

		// Flags with numeric values
		{"timeout flag", "nslookup -timeout=5 example.com", "nslookup -timeout=N <name>"},
		{"retry flag", "nslookup -retry=3 example.com", "nslookup -retry=N <name>"},
		{"port flag", "nslookup -port=5353 example.com", "nslookup -port=N <name>"},

		// Boolean flags
		{"debug flag", "nslookup -debug example.com", "nslookup -debug <name>"},
		{"vc flag", "nslookup -vc example.com", "nslookup -vc <name>"},

		// Combined
		{"vc with type", "nslookup -vc -type=ANY example.com", "nslookup -vc -type=ANY <name>"},
		{"debug with type and server", "nslookup -debug -type=MX example.com 8.8.8.8", "nslookup -debug -type=MX <name> <server>"},
		{"timeout and type", "nslookup -timeout=10 -type=NS example.com", "nslookup -timeout=N -type=NS <name>"},
		{"port type server", "nslookup -port=5353 -type=TXT example.com ns1.example.com", "nslookup -port=N -type=TXT <name> <server>"},
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
		a := Normalize("nslookup example.com")
		b := Normalize("nslookup google.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different servers collide", func(t *testing.T) {
		a := Normalize("nslookup example.com 8.8.8.8")
		b := Normalize("nslookup example.com 1.1.1.1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different timeouts collide", func(t *testing.T) {
		a := Normalize("nslookup -timeout=5 example.com")
		b := Normalize("nslookup -timeout=30 google.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("nslookup literal-arg")
		subshell := Normalize("nslookup $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
