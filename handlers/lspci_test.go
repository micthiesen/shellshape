package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLspci(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare", "lspci", "lspci"},
		{"verbose", "lspci -v", "lspci -v"},
		{"very verbose", "lspci -vv", "lspci -vv"},
		{"numeric ids", "lspci -nn", "lspci -nn"},
		{"kernel drivers", "lspci -k", "lspci -k"},
		{"tree", "lspci -t", "lspci -t"},
		{"slot filter", "lspci -s 00:1f.0", "lspci -s <val>"},
		{"device filter", "lspci -d 8086:a370", "lspci -d <val>"},
		{"kernel with slot", "lspci -k -s 01:00.0", "lspci -k -s <val>"},
		{"ids file", "lspci -i /usr/share/pci.ids", "lspci -i <path>"},
		{"proc dir", "lspci -p /proc/bus/pci", "lspci -p <path>"},
		{"dump file", "lspci -F /tmp/lspci.dump", "lspci -F <path>"},
		{"access method", "lspci -A intel-conf1", "lspci -A <val>"},
		{"combined flags", "lspci -v -nn -k", "lspci -v -nn -k"},
		{"with redirect", "lspci -v > /tmp/pci.txt", "lspci -v > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different slots collide", func(t *testing.T) {
		a := shellshape.Normalize("lspci -s 00:1f.0")
		b := shellshape.Normalize("lspci -s 01:00.0")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lspci -s 00:1f.0")
		subshell := shellshape.Normalize("lspci -s $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
