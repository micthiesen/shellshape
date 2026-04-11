package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLsof(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "lsof", "lsof"},
		{"file path", "lsof /var/log/syslog", "lsof <path>"},
		{"multiple files", "lsof /tmp/foo.sock /var/run/daemon.pid", "lsof <path>+"},
		{"terse flag", "lsof -t", "lsof -t"},

		// Flags with arguments
		{"network spec port", "lsof -i :8080", "lsof -i <net-spec>"},
		{"network spec tcp", "lsof -i TCP:*", "lsof -i <net-spec>"},
		{"pid", "lsof -p 1234", "lsof -p <pid>"},
		{"user", "lsof -u root", "lsof -u <user>"},
		{"command name", "lsof -c nginx", "lsof -c <name>"},
		{"fd numbers", "lsof -d 0-2", "lsof -d <fd>"},
		{"cache dir", "lsof -D /tmp/cache", "lsof -D <path>"},
		{"plus D directory", "lsof +D /var/log", "lsof +D <path>"},

		// Boolean -i (no value)
		{"bare -i", "lsof -i", "lsof -i"},
		{"bare -i with other flag", "lsof -i -P", "lsof -i -P"},

		// Combined flags
		{"terse with network", "lsof -t -i :80", "lsof -t -i <net-spec>"},
		{"pid and network", "lsof -p 5678 -a -i :443", "lsof -p <pid> -a -i <net-spec>"},
		{"user and command", "lsof -u www-data -c apache2", "lsof -u <user> -c <name>"},
		{"verbose network", "lsof -i -P -n", "lsof -i -P -n"},

		// Boolean flags preserved
		{"and flag", "lsof -a -u root -i", "lsof -a -u <user> -i"},
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
	t.Run("different PIDs collide", func(t *testing.T) {
		a := shellshape.Normalize("lsof -p 1234")
		b := shellshape.Normalize("lsof -p 9999")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different users collide", func(t *testing.T) {
		a := shellshape.Normalize("lsof -u root")
		b := shellshape.Normalize("lsof -u www-data")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different network specs collide", func(t *testing.T) {
		a := shellshape.Normalize("lsof -i :8080")
		b := shellshape.Normalize("lsof -i :3000")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lsof /var/log/syslog")
		subshell := shellshape.Normalize("lsof $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestLsofRedirects(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"with redirect", "lsof -i :8080 > /tmp/out.txt", "lsof -i <net-spec> > <path>"},
		{"with stderr redirect", "lsof -p 123 2>/dev/null", "lsof -p <pid> 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
