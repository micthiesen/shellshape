package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestIperf3(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Server mode
		{"server", "iperf3 -s", "iperf3 -s"},
		{"server with port", "iperf3 -s -p 5201", "iperf3 -s -p N"},
		{"server daemon", "iperf3 -s -D", "iperf3 -s -D"},

		// Client mode
		{"client basic", "iperf3 -c 192.168.1.100", "iperf3 -c <val>"},
		{"client hostname", "iperf3 -c server.example.com", "iperf3 -c <val>"},
		{"client with time", "iperf3 -c 10.0.0.1 -t 30", "iperf3 -c <val> -t N"},
		{"client with parallel", "iperf3 -c 10.0.0.1 -P 4", "iperf3 -c <val> -P N"},
		{"client with port", "iperf3 -c 10.0.0.1 -p 5201", "iperf3 -c <val> -p N"},
		{"client reverse", "iperf3 -c 10.0.0.1 -R", "iperf3 -c <val> -R"},
		{"client udp", "iperf3 -c 10.0.0.1 -u", "iperf3 -c <val> -u"},
		{"client bandwidth", "iperf3 -c 10.0.0.1 -b 100M", "iperf3 -c <val> -b <val>"},
		{"client bytes", "iperf3 -c 10.0.0.1 -n 1G", "iperf3 -c <val> -n N"},
		{"client bind", "iperf3 -c 10.0.0.1 --bind 192.168.1.50", "iperf3 -c <val> --bind <val>"},

		// Complex invocations
		{"client full", "iperf3 -c 10.0.0.1 -u -b 50M -t 60 -P 8 -R", "iperf3 -c <val> -u -b <val> -t N -P N -R"},
		{"client with interval", "iperf3 -c 10.0.0.1 -i 2", "iperf3 -c <val> -i N"},
		{"client with format", "iperf3 -c 10.0.0.1 -f m", "iperf3 -c <val> -f <val>"},

		// Boolean flags
		{"json output", "iperf3 -c 10.0.0.1 -J", "iperf3 -c <val> -J"},
		{"verbose", "iperf3 -c 10.0.0.1 -V", "iperf3 -c <val> -V"},
		{"version", "iperf3 --version", "iperf3 --version"},

		// Long flags
		{"long client", "iperf3 --client 10.0.0.1", "iperf3 --client <val>"},
		{"long server", "iperf3 --server", "iperf3 --server"},
		{"long port", "iperf3 -s --port 5201", "iperf3 -s --port N"},

		// Log file
		{"logfile", "iperf3 -s --logfile /var/log/iperf3.log", "iperf3 -s --logfile <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different hosts -> same shape
	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("iperf3 -c 192.168.1.100 -t 10")
		b := shellshape.Normalize("iperf3 -c server.example.com -t 30")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("iperf3 -c myhost")
		subshell := shellshape.Normalize("iperf3 -c $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
