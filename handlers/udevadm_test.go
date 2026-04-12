package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUdevadm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"control reload-rules", "udevadm control --reload-rules", "udevadm control --reload-rules"},
		{"control reload", "udevadm control --reload", "udevadm control --reload"},
		{"settle", "udevadm settle", "udevadm settle"},

		// Trigger subcommand
		{"trigger type", "udevadm trigger --type=subsystems", "udevadm trigger --type=subsystems"},
		{"trigger action add", "udevadm trigger --action add", "udevadm trigger --action add"},
		{"trigger action change", "udevadm trigger --action change", "udevadm trigger --action change"},

		// Info subcommand
		{"info query name", "udevadm info -q property -n /dev/sda", "udevadm info -q property -n <path>"},
		{"info path", "udevadm info -p /sys/class/net/eth0", "udevadm info -p <path>"},
		{"info long flags", "udevadm info --query=name --name=/dev/nvme0n1", "udevadm info --query=name --name=<path>"},

		// Monitor subcommand
		{"monitor property", "udevadm monitor --property", "udevadm monitor --property"},
		{"monitor kernel udev", "udevadm monitor --kernel --udev", "udevadm monitor --kernel --udev"},

		// Test and hwdb subcommands
		{"test device", "udevadm test /sys/class/net/eth0", "udevadm test <path>"},
		{"hwdb update", "udevadm hwdb --update", "udevadm hwdb --update"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different device paths collide", func(t *testing.T) {
		a := shellshape.Normalize("udevadm info -n /dev/sda")
		b := shellshape.Normalize("udevadm info -n /dev/nvme0n1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("udevadm info -n /dev/sda")
		subshell := shellshape.Normalize("udevadm info -n $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
