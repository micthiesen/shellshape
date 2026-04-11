package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple host", "mtr google.com", "mtr <host>"},
		{"simple ip", "mtr 8.8.8.8", "mtr <host>"},

		// Flags with numeric arguments
		{"report cycles", "mtr -c 10 google.com", "mtr -c N <host>"},
		{"interval", "mtr -i 0.5 10.0.0.1", "mtr -i N <host>"},
		{"packet size", "mtr -s 1500 example.com", "mtr -s N <host>"},
		{"max ttl", "mtr -m 20 example.com", "mtr -m N <host>"},
		{"first ttl", "mtr -f 5 example.com", "mtr -f N <host>"},
		{"port", "mtr -p 443 example.com", "mtr -p N <host>"},
		{"long report cycles", "mtr --report-cycles 10 google.com", "mtr --report-cycles N <host>"},
		{"long max ttl", "mtr --max-ttl 15 example.com", "mtr --max-ttl N <host>"},
		{"tos", "mtr -Q 16 example.com", "mtr -Q N <host>"},
		{"bitpattern", "mtr -B 255 example.com", "mtr -B N <host>"},
		{"timeout", "mtr -z 5 example.com", "mtr -z N <host>"},

		// Flags with value arguments
		{"source address", "mtr -a 192.168.1.100 example.com", "mtr -a <val> <host>"},
		{"interface", "mtr -I eth0 example.com", "mtr -I <val> <host>"},
		{"field order", "mtr -o LSR example.com", "mtr -o <val> <host>"},
		{"output file", "mtr -F /tmp/report.txt example.com", "mtr -F <val> <host>"},
		{"long interface", "mtr --interface eth0 example.com", "mtr --interface <val> <host>"},
		{"long address", "mtr --address 10.0.0.1 example.com", "mtr --address <val> <host>"},

		// Boolean flags
		{"report", "mtr -r google.com", "mtr -r <host>"},
		{"report wide", "mtr -w google.com", "mtr -w <host>"},
		{"no dns", "mtr -n google.com", "mtr -n <host>"},
		{"show ips", "mtr --show-ips example.com", "mtr --show-ips <host>"},
		{"tcp mode", "mtr --tcp example.com", "mtr --tcp <host>"},
		{"udp mode", "mtr --udp example.com", "mtr --udp <host>"},
		{"json output", "mtr --json example.com", "mtr --json <host>"},
		{"aslookup", "mtr --aslookup example.com", "mtr --aslookup <host>"},
		{"ipv4", "mtr -4 example.com", "mtr -4 <host>"},
		{"ipv6", "mtr -6 example.com", "mtr -6 <host>"},
		{"long report", "mtr --report example.com", "mtr --report <host>"},

		// Combinations
		{"report with cycles", "mtr -r -c 10 google.com", "mtr -r -c N <host>"},
		{"tcp with port", "mtr --tcp -p 443 example.com", "mtr --tcp -p N <host>"},
		{"no dns report wide", "mtr -n -w example.com", "mtr -n -w <host>"},
		{"multiple numeric", "mtr -c 5 -i 0.5 -s 64 example.com", "mtr -c N -i N -s N <host>"},
		{"full combo", "mtr -4 -r -c 10 -i 0.2 --tcp -p 80 example.com", "mtr -4 -r -c N -i N --tcp -p N <host>"},

		// Redirect
		{"redirect", "mtr -r -c 10 example.com > report.txt", "mtr -r -c N <host> > <path>"},
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
		a := shellshape.Normalize("mtr -c 10 google.com")
		b := shellshape.Normalize("mtr -c 10 yahoo.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different cycle counts collide", func(t *testing.T) {
		a := shellshape.Normalize("mtr -c 5 example.com")
		b := shellshape.Normalize("mtr -c 100 example.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mtr example.com")
		subshell := shellshape.Normalize("mtr $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
