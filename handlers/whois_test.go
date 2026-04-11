package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWhois(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple domain", "whois example.com", "whois <name>"},
		{"ip address", "whois 8.8.8.8", "whois <name>"},

		// Boolean flags
		{"arin flag", "whois -a example.com", "whois -a <name>"},
		{"recursive flag", "whois -R example.com", "whois -R <name>"},
		{"quick flag", "whois -Q example.com", "whois -Q <name>"},
		{"bundled flags", "whois -aR example.com", "whois -aR <name>"},

		// Flags with arguments
		{"custom host", "whois -h whois.verisign-grs.com example.com", "whois -h <host> <name>"},
		{"custom port", "whois -p 4321 example.com", "whois -p N <name>"},
		{"tld flag", "whois -c RU CONTACT-ID", "whois -c <val> <name>"},
		{"host and port", "whois -h whois.example.com -p 4321 query-data", "whois -h <host> -p N <name>"},

		// Multiple positionals collapse
		{"multiple names", "whois example.com example.org example.net", "whois <name>"},

		// Double dash stops flag parsing
		{"double dash", "whois -r -- example.com", "whois -r -- <name>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different domains produce the same shape
	t.Run("different domains collide", func(t *testing.T) {
		a := shellshape.Normalize("whois google.com")
		b := shellshape.Normalize("whois example.org")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("whois -h whois.arin.net example.com")
		b := shellshape.Normalize("whois -h whois.ripe.net example.org")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("whois example.com")
		subshell := shellshape.Normalize("whois $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
