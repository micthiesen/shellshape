package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSmartctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic info modes
		{"health", "smartctl -H /dev/sda", "smartctl -H <path>"},
		{"info", "smartctl -i /dev/sda", "smartctl -i <path>"},
		{"all", "smartctl -a /dev/sda", "smartctl -a <path>"},
		{"xall", "smartctl -x /dev/nvme0", "smartctl -x <path>"},
		{"attributes", "smartctl -A /dev/sda", "smartctl -A <path>"},
		{"capabilities", "smartctl -c /dev/sda", "smartctl -c <path>"},
		{"long flags", "smartctl --all /dev/sda", "smartctl --all <path>"},

		// Test types (structural)
		{"test short", "smartctl -t short /dev/sda", "smartctl -t short <path>"},
		{"test long", "smartctl -t long /dev/sda", "smartctl -t long <path>"},
		{"test conveyance", "smartctl --test=conveyance /dev/sda", "smartctl --test=conveyance <path>"},

		// Device type (structural)
		{"device type", "smartctl -d sat -a /dev/sda", "smartctl -d sat -a <path>"},
		{"device type nvme", "smartctl -d nvme -i /dev/nvme0", "smartctl -d nvme -i <path>"},
		{"device type long", "smartctl --device=ata -a /dev/sda", "smartctl --device=ata -a <path>"},

		// Log type (structural)
		{"log error", "smartctl -l error /dev/sda", "smartctl -l error <path>"},
		{"log selftest", "smartctl -l selftest /dev/sda", "smartctl -l selftest <path>"},
		{"log long", "smartctl --log=error /dev/sda", "smartctl --log=error <path>"},

		// Scan
		{"scan", "smartctl --scan", "smartctl --scan"},
		{"scan-open", "smartctl --scan-open", "smartctl --scan-open"},

		// Abort
		{"abort", "smartctl -X /dev/sda", "smartctl -X <path>"},

		// Nocheck (structural power mode)
		{"nocheck", "smartctl -n standby -a /dev/sda", "smartctl -n standby -a <path>"},

		// Combined flags
		{"combined", "smartctl -d sat -H -l error /dev/sdb", "smartctl -d sat -H -l error <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("smartctl -a /dev/sda")
		b := shellshape.Normalize("smartctl -a /dev/nvme0n1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("smartctl -a /dev/sda")
		subshell := shellshape.Normalize("smartctl -a $(echo /dev/sda)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
