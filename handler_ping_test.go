package shellshape

import "testing"

func TestPing(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple host", "ping google.com", "ping <host>"},
		{"simple ip", "ping 192.168.1.1", "ping <host>"},

		// Flags with numeric arguments
		{"count", "ping -c 5 google.com", "ping -c N <host>"},
		{"interval", "ping -i 0.5 10.0.0.1", "ping -i N <host>"},
		{"timeout", "ping -t 30 example.com", "ping -t N <host>"},
		{"packet size", "ping -s 1024 example.com", "ping -s N <host>"},
		{"waittime", "ping -W 1000 8.8.8.8", "ping -W N <host>"},
		{"ttl", "ping -m 64 example.com", "ping -m N <host>"},
		{"multiple numeric", "ping -c 10 -i 0.2 -s 512 host.local", "ping -c N -i N -s N <host>"},

		// Flags with non-numeric arguments
		{"source addr", "ping -S 192.168.1.100 example.com", "ping -S <val> <host>"},
		{"interface", "ping -I eth0 224.0.0.1", "ping -I <val> <host>"},
		{"bind interface", "ping -b en0 example.com", "ping -b <val> <host>"},
		{"mode flag", "ping -M mask example.com", "ping -M <val> <host>"},
		{"pattern", "ping -p ff example.com", "ping -p <val> <host>"},

		// Boolean flags
		{"quiet", "ping -q 8.8.8.8", "ping -q <host>"},
		{"numeric and quiet", "ping -n -q 8.8.8.8", "ping -n -q <host>"},
		{"flood", "ping -f -c 100 host.example.com", "ping -f -c N <host>"},
		{"verbose", "ping -v example.com", "ping -v <host>"},

		// Redirect
		{"redirect", "ping -c 3 example.com > out.txt", "ping -c N <host> > <path>"},
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
	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("ping -c 5 google.com")
		b := Normalize("ping -c 5 yahoo.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different counts collide", func(t *testing.T) {
		a := Normalize("ping -c 5 example.com")
		b := Normalize("ping -c 100 example.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ping google.com")
		subshell := Normalize("ping $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
