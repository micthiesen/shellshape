package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSudo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple command", "sudo ls /home/user", "sudo ls <path>"},
		{"no args", "sudo -i", "sudo -i"},
		{"just list", "sudo -l", "sudo -l"},
		{"reset timestamp", "sudo -K", "sudo -K"},

		// Flags with arguments
		{"user flag", "sudo -u www vi /var/www/index.html", "sudo -u <val> vi <path>"},
		{"group flag", "sudo -g wheel id", "sudo -g <val> id"},
		{"prompt flag", "sudo -p 'Enter password:' apt update", "sudo -p <str> apt update"},
		{"close-from", "sudo -C 4 ls", "sudo -C N ls"},

		// Boolean flags
		{"preserve env", "sudo -E make build", "sudo -E make build"},
		{"non-interactive", "sudo -n apt install vim", "sudo -n apt install vim"},
		{"background", "sudo -b /opt/server start", "sudo -b <path> start"},

		// Multiple flags
		{"user and env", "sudo -u deploy -E /opt/app/run.sh", "sudo -u <val> -E <path>"},
		{"non-interactive user", "sudo -n -u postgres psql", "sudo -n -u <val> psql"},

		// Inner command with paths and values
		{"inner paths", "sudo cp /etc/hosts /tmp/hosts.bak", "sudo cp <path>+"},
		{"inner flags", "sudo rm -rf /tmp/cache", "sudo rm -rf <path>"},

		// Long flags
		{"long user", "sudo --user=www ls", "sudo --user=<val> ls"},
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
	t.Run("different users collide", func(t *testing.T) {
		a := shellshape.Normalize("sudo -u alice ls /home")
		b := shellshape.Normalize("sudo -u bob ls /home")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different inner paths collide", func(t *testing.T) {
		a := shellshape.Normalize("sudo cat /etc/shadow")
		b := shellshape.Normalize("sudo cat /etc/passwd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sudo literal-arg")
		subshell := shellshape.Normalize("sudo $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
