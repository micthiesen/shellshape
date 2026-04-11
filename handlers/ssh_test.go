package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSsh(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple host", "ssh myserver", "ssh <host>"},
		{"user@host", "ssh user@example.com", "ssh <host>"},
		{"host with domain", "ssh web01.prod.example.com", "ssh <host>"},

		// Flags with arguments
		{"identity file", "ssh -i ~/.ssh/key.pem myserver", "ssh -i <path> <host>"},
		{"port", "ssh -p 2222 myserver", "ssh -p N <host>"},
		{"config file", "ssh -F /etc/ssh/config myserver", "ssh -F <path> <host>"},
		{"jump host", "ssh -J jumphost targethost", "ssh -J <jump-host> <host>"},
		{"login name", "ssh -l admin myserver", "ssh -l <str> <host>"},
		{"option", "ssh -o StrictHostKeyChecking=no myserver", "ssh -o <str> <host>"},
		{"local forward", "ssh -L 8080:localhost:3000 myserver", "ssh -L <str> <host>"},
		{"remote forward", "ssh -R 9090:localhost:80 myserver", "ssh -R <str> <host>"},
		{"dynamic forward", "ssh -D 1080 myserver", "ssh -D <str> <host>"},
		{"log file", "ssh -E /tmp/ssh.log myserver", "ssh -E <path> <host>"},

		// Boolean flags
		{"verbose", "ssh -v myserver", "ssh -v <host>"},
		{"very verbose", "ssh -vvv myserver", "ssh -vvv <host>"},
		{"force tty", "ssh -t myserver", "ssh -t <host>"},
		{"compression", "ssh -C myserver", "ssh -C <host>"},
		{"agent forwarding", "ssh -A myserver", "ssh -A <host>"},
		{"no remote cmd", "ssh -N myserver", "ssh -N <host>"},

		// Remote command
		{"remote command", "ssh myserver ls -la /tmp", "ssh <host> <cmd>+"},
		{"remote command single", "ssh myserver uptime", "ssh <host> <cmd>"},
		{"remote tty command", "ssh -t myserver sudo reboot", "ssh -t <host> <cmd>+"},

		// Combined flags
		{"port and identity", "ssh -p 22 -i ~/.ssh/id_rsa myserver", "ssh -p N -i <path> <host>"},
		{"forward and verbose", "ssh -v -L 8080:localhost:80 myserver", "ssh -v -L <str> <host>"},
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
		a := shellshape.Normalize("ssh production-server")
		b := shellshape.Normalize("ssh staging-server")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different remote commands collide", func(t *testing.T) {
		a := shellshape.Normalize("ssh myserver cat /etc/hosts")
		b := shellshape.Normalize("ssh otherserver tail -f /var/log/syslog")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := shellshape.Normalize("ssh -p 22 myserver")
		b := shellshape.Normalize("ssh -p 2222 otherserver")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ssh literal-host")
		subshell := shellshape.Normalize("ssh $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
