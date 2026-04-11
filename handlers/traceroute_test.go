package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTraceroute(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple host", "traceroute google.com", "traceroute <host>"},
		{"simple ip", "traceroute 8.8.8.8", "traceroute <host>"},
		{"with packet size", "traceroute google.com 60", "traceroute <host> N"},

		// Numeric flags
		{"max hops", "traceroute -m 30 example.com", "traceroute -m N <host>"},
		{"queries per hop", "traceroute -q 1 example.com", "traceroute -q N <host>"},
		{"wait time", "traceroute -w 3 example.com", "traceroute -w N <host>"},
		{"port", "traceroute -p 8080 example.com", "traceroute -p N <host>"},
		{"first ttl", "traceroute -f 5 example.com", "traceroute -f N <host>"},
		{"first ttl M", "traceroute -M 2 example.com", "traceroute -M N <host>"},
		{"tos", "traceroute -t 16 example.com", "traceroute -t N <host>"},
		{"pause", "traceroute -z 500 example.com", "traceroute -z N <host>"},

		// Value flags
		{"source addr", "traceroute -s 10.0.0.1 example.com", "traceroute -s <val> <host>"},
		{"interface", "traceroute -i eth0 example.com", "traceroute -i <val> <host>"},
		{"gateway", "traceroute -g 192.168.1.1 example.com", "traceroute -g <val> <host>"},
		{"as server", "traceroute -A whois.example.com example.com", "traceroute -A <val> <host>"},
		{"protocol", "traceroute -P icmp example.com", "traceroute -P <val> <host>"},

		// Boolean flags
		{"icmp", "traceroute -I example.com", "traceroute -I <host>"},
		{"tcp", "traceroute -T example.com", "traceroute -T <host>"},
		{"no dns", "traceroute -n example.com", "traceroute -n <host>"},
		{"ipv4", "traceroute -4 example.com", "traceroute -4 <host>"},
		{"ipv6", "traceroute -6 example.com", "traceroute -6 <host>"},

		// Combinations
		{"icmp no dns max hops", "traceroute -I -n -m 20 example.com", "traceroute -I -n -m N <host>"},
		{"multiple numeric", "traceroute -q 1 -w 3 -m 15 host.local", "traceroute -q N -w N -m N <host>"},
		{"value and numeric", "traceroute -i eth0 -m 30 -q 2 example.com 100", "traceroute -i <val> -m N -q N <host> N"},

		// Alias
		{"tracepath alias", "tracepath example.com", "tracepath <host>"},
		{"tracepath with flag", "tracepath -n example.com", "tracepath -n <host>"},

		// Redirect
		{"redirect", "traceroute example.com > out.txt", "traceroute <host> > <path>"},
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
	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("traceroute -m 30 google.com")
		b := shellshape.Normalize("traceroute -m 30 yahoo.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different packet sizes collide", func(t *testing.T) {
		a := shellshape.Normalize("traceroute example.com 60")
		b := shellshape.Normalize("traceroute example.com 1500")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("traceroute example.com")
		subshell := shellshape.Normalize("traceroute $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
