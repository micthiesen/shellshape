package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestHwinfo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare", "hwinfo", "hwinfo"},
		{"short cpu", "hwinfo --short --cpu", "hwinfo --short --cpu"},
		{"short gfxcard", "hwinfo --short --gfxcard", "hwinfo --short --gfxcard"},
		{"all", "hwinfo --all", "hwinfo --all"},
		{"short disk", "hwinfo --short --disk", "hwinfo --short --disk"},
		{"memory", "hwinfo --short --memory", "hwinfo --short --memory"},
		{"log file", "hwinfo --log /tmp/hw.log", "hwinfo --log <path>"},
		{"only device", "hwinfo --only /dev/sda", "hwinfo --only <path>"},
		{"log fused", "hwinfo --all --log=/tmp/hw.log", "hwinfo --all --log=<path>"},
		{"only fused", "hwinfo --only=/dev/nvme0n1", "hwinfo --only=<path>"},
		{"multiple probes", "hwinfo --short --netcard --wlan", "hwinfo --short --netcard --wlan"},
		{"with redirect", "hwinfo --short --cpu > /tmp/cpu.txt", "hwinfo --short --cpu > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different log paths collide", func(t *testing.T) {
		a := shellshape.Normalize("hwinfo --log /tmp/a.log")
		b := shellshape.Normalize("hwinfo --log /var/log/hw.log")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("hwinfo --short --cpu")
		subshell := shellshape.Normalize("hwinfo $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
