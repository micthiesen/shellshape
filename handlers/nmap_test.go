package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNmap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple host", "nmap 192.168.1.1", "nmap <host>"},
		{"hostname", "nmap example.com", "nmap <host>"},
		{"cidr range", "nmap 10.0.0.0/24", "nmap <host>"},
		{"multiple targets", "nmap 10.0.0.1 10.0.0.2 10.0.0.3", "nmap <host>+"},

		// Scan types (boolean flags)
		{"syn scan", "nmap -sS 192.168.1.1", "nmap -sS <host>"},
		{"aggressive scan", "nmap -A -T4 192.168.1.1", "nmap -A -T4 <host>"},
		{"OS detection", "nmap -O --osscan-guess 10.0.0.1", "nmap -O --osscan-guess <host>"},
		{"service version", "nmap -sV 192.168.1.1", "nmap -sV <host>"},

		// Port specification
		{"port flag", "nmap -p 22,80,443 192.168.1.1", "nmap -p <ports> <host>"},
		{"port range", "nmap -p 1-65535 10.0.0.1", "nmap -p <ports> <host>"},
		{"top ports", "nmap --top-ports 100 192.168.1.1", "nmap --top-ports N <host>"},

		// Output flags
		{"output normal", "nmap -oN output.txt 192.168.1.1", "nmap -oN <path> <host>"},
		{"output xml", "nmap -oX scan.xml 10.0.0.1", "nmap -oX <path> <host>"},
		{"output all", "nmap -oA scanresults 192.168.1.1", "nmap -oA <path> <host>"},

		// Input list
		{"input list", "nmap -iL targets.txt", "nmap -iL <path>"},

		// Script flags
		{"script", "nmap --script vuln 192.168.1.1", "nmap --script <script> <host>"},
		{"script args", "nmap --script smb-enum-shares --script-args user=admin 10.0.0.1", "nmap --script <script> --script-args <val> <host>"},

		// Numeric flags
		{"min rate", "nmap --min-rate 5000 10.0.0.1", "nmap --min-rate N <host>"},
		{"max retries", "nmap --max-retries 3 192.168.1.1", "nmap --max-retries N <host>"},
		{"source port", "nmap -g 53 192.168.1.1", "nmap -g N <host>"},
		{"data length", "nmap --data-length 25 10.0.0.1", "nmap --data-length N <host>"},

		// Value flags
		{"decoy", "nmap -D RND:10 192.168.1.1", "nmap -D <val> <host>"},
		{"interface", "nmap -e eth0 192.168.1.1", "nmap -e <val> <host>"},
		{"spoof mac", "nmap --spoof-mac 00:11:22:33:44:55 10.0.0.1", "nmap --spoof-mac <val> <host>"},
		{"exclude", "nmap --exclude 10.0.0.5 10.0.0.0/24", "nmap --exclude <host>+"},

		// Source IP
		{"source ip", "nmap -S 10.0.0.99 192.168.1.1", "nmap -S <host>+"},

		// Complex real-world
		{"complex scan", "nmap -sS -p 1-65535 --min-rate 5000 -T4 -n 10.0.0.1", "nmap -sS -p <ports> --min-rate N -T4 -n <host>"},
		{"heartbleed", "nmap -T5 --min-parallelism 50 -n --script ssl-heartbleed -p 443 127.0.0.1", "nmap -T5 --min-parallelism N -n --script <script> -p <ports> <host>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different targets collide", func(t *testing.T) {
		a := shellshape.Normalize("nmap -sS -p 80 192.168.1.1")
		b := shellshape.Normalize("nmap -sS -p 80 10.0.0.5")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := shellshape.Normalize("nmap -p 22,80 192.168.1.1")
		b := shellshape.Normalize("nmap -p 443,8080 10.0.0.1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nmap 192.168.1.1")
		subshell := shellshape.Normalize("nmap $(cat targets.txt)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
