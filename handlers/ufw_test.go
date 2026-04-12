package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUfw(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"enable", "ufw enable", "ufw enable"},
		{"disable", "ufw disable", "ufw disable"},
		{"status", "ufw status", "ufw status"},
		{"status numbered", "ufw status numbered", "ufw status numbered"},
		{"status verbose", "ufw status verbose", "ufw status verbose"},
		{"reset", "ufw reset", "ufw reset"},
		{"reload", "ufw reload", "ufw reload"},

		// Simple allow/deny with port
		{"allow port", "ufw allow 22", "ufw allow N"},
		{"deny port", "ufw deny 80", "ufw deny N"},
		{"reject port", "ufw reject 443", "ufw reject N"},
		{"limit port", "ufw limit 22", "ufw limit N"},
		{"allow port proto", "ufw allow 22/tcp", "ufw allow <val>"},
		{"allow port range", "ufw allow 8000:9000/tcp", "ufw allow <val>"},

		// Allow with from/to syntax
		{"allow from ip", "ufw allow from 192.168.1.0/24", "ufw allow from <val>"},
		{"allow from any to any port", "ufw allow from any to any port 22", "ufw allow from any to any port N"},
		{"allow proto from to port", "ufw allow proto tcp from 192.168.0.4 to any port 22", "ufw allow proto tcp from <val> to any port N"},
		{"deny from ip", "ufw deny from 10.0.0.0/8", "ufw deny from <val>"},
		{"deny proto udp range", "ufw deny proto udp from any to any port 8412:8500", "ufw deny proto udp from any to any port <val>"},

		// Delete rules
		{"delete by number", "ufw delete 3", "ufw delete N"},
		{"delete rule", "ufw delete allow 22", "ufw delete allow N"},
		{"delete deny rule", "ufw delete deny from 10.0.0.0/8", "ufw delete deny from <val>"},

		// App rules
		{"app list", "ufw app list", "ufw app list"},
		{"app info", "ufw app info OpenSSH", "ufw app info <val>"},
		{"allow app", "ufw allow OpenSSH", "ufw allow <val>"},

		// Comment
		{"allow with comment", "ufw allow 5432 comment Service", "ufw allow N comment <val>"},

		// Direction
		{"allow in on", "ufw allow in on eth0", "ufw allow in on <val>"},
		{"allow out", "ufw allow out 80", "ufw allow out N"},

		// Flags
		{"dry-run", "ufw --dry-run allow 22", "ufw --dry-run allow N"},
		{"force enable", "ufw --force enable", "ufw --force enable"},

		// Logging
		{"logging on", "ufw logging on", "ufw logging on"},
		{"logging medium", "ufw logging medium", "ufw logging <val>"},

		// Default
		{"default deny", "ufw default deny incoming", "ufw default deny incoming"},
		{"default allow", "ufw default allow outgoing", "ufw default allow outgoing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different ports -> same shape
	t.Run("different ports collide", func(t *testing.T) {
		a := shellshape.Normalize("ufw allow 22")
		b := shellshape.Normalize("ufw allow 443")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different IPs collide", func(t *testing.T) {
		a := shellshape.Normalize("ufw allow from 10.0.0.0/8")
		b := shellshape.Normalize("ufw allow from 192.168.1.0/24")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ufw allow 22")
		subshell := shellshape.Normalize("ufw allow $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
