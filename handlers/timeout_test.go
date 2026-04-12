package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTimeout(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple duration and command", "timeout 5s sleep 10", "timeout N sleep N"},
		{"numeric duration", "timeout 30 ls /home/user", "timeout N ls <path>"},
		{"duration with suffix", "timeout 2m ./long-script.sh", "timeout N <path>"},

		// Flags with arguments
		{"signal flag short", "timeout -s INT 5s sleep 10", "timeout -s INT N sleep N"},
		{"signal flag long", "timeout --signal KILL 10 /usr/bin/app", "timeout --signal KILL N <path>"},
		{"kill-after short", "timeout -k 30s 5m ./task.sh", "timeout -k N N <path>"},
		{"kill-after long", "timeout --kill-after 10s 60s make build", "timeout --kill-after N N make build"},

		// Boolean flags
		{"foreground", "timeout --foreground 10 /usr/bin/server", "timeout --foreground N <path>"},
		{"preserve-status", "timeout --preserve-status 5s curl http://example.com/api", "timeout --preserve-status N curl <http-uri>"},
		{"verbose", "timeout -v 30s dd if=/dev/zero of=/tmp/out.bin", "timeout -v N dd <path>+"},

		// Multiple flags combined
		{"signal and kill-after", "timeout -s TERM -k 5s 30s /opt/app/run.sh --config /etc/app.conf", "timeout -s TERM -k N N <path> --config <path>"},
		{"foreground and verbose", "timeout --foreground -v 60 python3 server.py", "timeout --foreground -v N python3 <path>"},

		// Inner command with flags
		{"inner command flags", "timeout 10s grep -r pattern /var/log", "timeout N grep -r pattern <path>"},
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
	t.Run("different durations collide", func(t *testing.T) {
		a := shellshape.Normalize("timeout 5s sleep 10")
		b := shellshape.Normalize("timeout 30m sleep 10")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different inner paths collide", func(t *testing.T) {
		a := shellshape.Normalize("timeout 10s /usr/bin/app1")
		b := shellshape.Normalize("timeout 10s /usr/bin/app2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("timeout 5s literal-arg")
		subshell := shellshape.Normalize("timeout 5s $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
