package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNohup(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple command", "nohup ls", "nohup ls"},
		{"command with flag", "nohup ls -l", "nohup ls -l"},
		{"script path", "nohup ./deploy.sh", "nohup <path>"},
		{"absolute script path", "nohup /usr/local/bin/server", "nohup <path>"},

		// Command with arguments
		{"command with path arg", "nohup tail -f /var/log/syslog", "nohup tail -f <path>"},
		{"command with multiple args", "nohup python3 app.py --port 8080", "nohup python3 <path> --port N"},
		{"java jar", "nohup java -jar server.jar", "nohup java -jar <dotted-id>"},
		{"sleep with number", "nohup sleep 3600", "nohup sleep N"},

		// Redirects
		{"with output redirect", "nohup ./run.sh > /tmp/out.log", "nohup <path> > <path>"},
		{"with stderr redirect", "nohup ./run.sh 2> /tmp/err.log", "nohup <path> 2> <path>"},
		{"with both redirects", "nohup ./run.sh > /tmp/out.log 2>&1", "nohup <path> > <path> 2>&1"},

		// Flags on the wrapped command
		{"grep with flags", "nohup grep -r pattern /var/log", "nohup grep -r pattern <path>"},
		{"curl with flag", "nohup curl -O https://example.com/file.tar.gz", "nohup curl -O <https-uri>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests: different data values -> same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("nohup /opt/app/server")
		b := shellshape.Normalize("nohup /usr/local/bin/daemon")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different script names collide", func(t *testing.T) {
		a := shellshape.Normalize("nohup ./start.sh")
		b := shellshape.Normalize("nohup ./deploy.sh")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nohup literal-arg")
		subshell := shellshape.Normalize("nohup $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
