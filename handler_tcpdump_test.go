package shellshape

import "testing"

func TestTcpdump(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"interface only", "tcpdump -i eth0", "tcpdump -i <val>"},
		{"list interfaces", "tcpdump -D", "tcpdump -D"},
		{"simple host filter", "tcpdump host 192.168.1.1", "tcpdump <filter>+"},
		{"simple port filter", "tcpdump port 80", "tcpdump <filter>+"},
		{"protocol filter", "tcpdump tcp", "tcpdump <filter>"},

		// Flags with arguments
		{"write file", "tcpdump -w capture.pcap", "tcpdump -w <path>"},
		{"read file", "tcpdump -r capture.pcap", "tcpdump -r <path>"},
		{"count", "tcpdump -c 10", "tcpdump -c N"},
		{"snaplen", "tcpdump -s 1500", "tcpdump -s N"},
		{"interface and filter", "tcpdump -i eth0 port 80", "tcpdump -i <val> <filter>+"},
		{"count and interface", "tcpdump -c 10 -i eth0", "tcpdump -c N -i <val>"},
		{"rotate seconds", "tcpdump -G 3600 -w trace.pcap", "tcpdump -G N -w <path>"},
		{"file count", "tcpdump -W 5 -w trace.pcap", "tcpdump -W N -w <path>"},
		{"buffer size", "tcpdump -B 4096", "tcpdump -B N"},
		{"filter file", "tcpdump -F /tmp/filter.bpf", "tcpdump -F <path>"},
		{"file list", "tcpdump -V /tmp/files.txt", "tcpdump -V <path>"},

		// Boolean flags
		{"verbose", "tcpdump -v", "tcpdump -v"},
		{"no resolve", "tcpdump -nn", "tcpdump -nn"},
		{"extra verbose", "tcpdump -vvv", "tcpdump -vvv"},
		{"timestamp", "tcpdump -tttt", "tcpdump -tttt"},

		// Combined flags and filter
		{"verbose interface filter", "tcpdump -nn -i any host 10.0.0.1 and port 80", "tcpdump -nn -i <val> <filter>+"},
		{"complex filter", "tcpdump -i eth0 src 192.168.1.1 and dst 10.0.0.2 and dst port 80", "tcpdump -i <val> <filter>+"},

		// Quoted filter expression (single token)
		{"quoted filter", "tcpdump 'tcp port 80 and host 10.0.0.1'", "tcpdump <filter>"},

		// Write with filter
		{"write with filter", "tcpdump -w dump.pcap port not 22", "tcpdump -w <path> <filter>+"},

		// Direction flag
		{"direction", "tcpdump -Q in", "tcpdump -Q <val>"},

		// Link type
		{"linktype", "tcpdump -y EN10MB", "tcpdump -y <val>"},

		// Packet type
		{"packet type", "tcpdump -T rpc", "tcpdump -T <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("tcpdump -i eth0 host 192.168.1.1")
		b := Normalize("tcpdump -i eth0 host 10.0.0.1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := Normalize("tcpdump port 80")
		b := Normalize("tcpdump port 443")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pcap files collide", func(t *testing.T) {
		a := Normalize("tcpdump -w /tmp/a.pcap")
		b := Normalize("tcpdump -w /tmp/b.pcap")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tcpdump -i eth0")
		subshell := Normalize("tcpdump -i $(get-interface)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
