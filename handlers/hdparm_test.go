package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestHdparm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic info modes (boolean flags)
		{"info", "hdparm -I /dev/sda", "hdparm -I <path>"},
		{"timing", "hdparm -t /dev/sda", "hdparm -t <path>"},
		{"timing cache", "hdparm -T /dev/sda", "hdparm -T <path>"},
		{"combined timing", "hdparm -tT /dev/sda", "hdparm -tT <path>"},
		{"power mode", "hdparm -C /dev/sda", "hdparm -C <path>"},
		{"geometry", "hdparm -g /dev/sda", "hdparm -g <path>"},

		// Numeric flags with separate values
		{"apm set", "hdparm -B 128 /dev/sda", "hdparm -B N <path>"},
		{"standby timeout", "hdparm -S 242 /dev/sda", "hdparm -S N <path>"},
		{"acoustic", "hdparm -M 128 /dev/sda", "hdparm -M N <path>"},
		{"readahead", "hdparm -a 256 /dev/sda", "hdparm -a N <path>"},
		{"multcount", "hdparm -m 16 /dev/sda", "hdparm -m N <path>"},
		{"write cache", "hdparm -W 1 /dev/sda", "hdparm -W N <path>"},

		// Combined boolean and numeric
		{"apm get", "hdparm -B /dev/sda", "hdparm -B <path>"},
		{"multiple flags", "hdparm -I -B 128 /dev/sda", "hdparm -I -B N <path>"},

		// Power modes (boolean)
		{"standby", "hdparm -y /dev/sda", "hdparm -y <path>"},
		{"sleep", "hdparm -Y /dev/sda", "hdparm -Y <path>"},

		// Security
		{"security set pass", "hdparm --security-set-pass secret /dev/sda", "hdparm --security-set-pass <val> <path>"},
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
		a := shellshape.Normalize("hdparm -I /dev/sda")
		b := shellshape.Normalize("hdparm -I /dev/sdb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric values collide", func(t *testing.T) {
		a := shellshape.Normalize("hdparm -B 128 /dev/sda")
		b := shellshape.Normalize("hdparm -B 254 /dev/sdb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("hdparm -I /dev/sda")
		subshell := shellshape.Normalize("hdparm -I $(echo /dev/sda)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
