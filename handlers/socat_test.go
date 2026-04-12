package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSocat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"stdin to tcp", "socat - TCP:localhost:8080", "socat <val>+"},
		{"tcp listen and forward", "socat TCP-LISTEN:8080,fork TCP:remote:80", "socat <val>+"},
		{"unix socket", "socat UNIX-LISTEN:/tmp/sock STDIN", "socat <val>+"},
		{"stdin stdout", "socat STDIN STDOUT", "socat <val>+"},
		{"exec address", "socat TCP-LISTEN:1234 EXEC:/bin/cat", "socat <val>+"},

		// Boolean flags
		{"debug flag", "socat -d TCP-LISTEN:80 TCP:host:80", "socat -d <val>+"},
		{"double debug", "socat -d -d TCP-LISTEN:80 TCP:host:80", "socat -d -d <val>+"},
		{"verbose flag", "socat -v TCP-LISTEN:80 TCP:host:80", "socat -v <val>+"},
		{"hex dump", "socat -x STDIN STDOUT", "socat -x <val>+"},
		{"unidirectional", "socat -u TCP:host:80 STDOUT", "socat -u <val>+"},

		// Numeric flags
		{"buffer size", "socat -b 4096 TCP:host:80 STDOUT", "socat -b N <val>+"},
		{"timeout", "socat -T 5 TCP:host:80 STDOUT", "socat -T N <val>+"},

		// Combined
		{"typical debug forward", "socat -d -d -v TCP-LISTEN:8080,fork,reuseaddr TCP:localhost:80",
			"socat -d -d -v <val>+"},

		// Redirect
		{"with redirect", "socat -u STDIN STDOUT 2>/dev/null", "socat -u <val>+ 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different addresses collide", func(t *testing.T) {
		a := shellshape.Normalize("socat TCP-LISTEN:8080 TCP:host1:80")
		b := shellshape.Normalize("socat TCP-LISTEN:9090 TCP:host2:443")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("socat STDIN STDOUT")
		subshell := shellshape.Normalize("socat $(dangerous-command) STDOUT")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
