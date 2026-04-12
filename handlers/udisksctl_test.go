package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUdisksctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"status", "udisksctl status", "udisksctl status"},
		{"monitor", "udisksctl monitor", "udisksctl monitor"},
		{"help", "udisksctl help", "udisksctl help"},
		{"dump", "udisksctl dump", "udisksctl dump"},

		// Mount/unmount with block device
		{"mount short", "udisksctl mount -b /dev/sda1", "udisksctl mount -b <path>"},
		{"mount long", "udisksctl mount --block-device /dev/sda1", "udisksctl mount --block-device <path>"},
		{"unmount short", "udisksctl unmount -b /dev/sdb1", "udisksctl unmount -b <path>"},
		{"mount with options", "udisksctl mount -b /dev/sda1 -t ext4 -o noatime", "udisksctl mount -b <path> -t <val> -o <val>"},
		{"mount no-user-interaction", "udisksctl mount --no-user-interaction -b /dev/sde1", "udisksctl mount --no-user-interaction -b <path>"},

		// Object path
		{"info object-path", "udisksctl info -p /org/freedesktop/UDisks2/drives/foo", "udisksctl info -p <path>"},
		{"info long object-path", "udisksctl info --object-path /org/freedesktop/UDisks2/block_devices/sda1", "udisksctl info --object-path <path>"},

		// Lock/unlock
		{"lock", "udisksctl lock -b /dev/sda2", "udisksctl lock -b <path>"},
		{"unlock", "udisksctl unlock -b /dev/sda2", "udisksctl unlock -b <path>"},

		// Loop setup/delete
		{"loop-setup", "udisksctl loop-setup -f /home/user/disk.img", "udisksctl loop-setup -f <path>"},
		{"loop-setup read-only", "udisksctl loop-setup --read-only -f /tmp/image.iso", "udisksctl loop-setup --read-only -f <path>"},
		{"loop-delete", "udisksctl loop-delete -b /dev/loop0", "udisksctl loop-delete -b <path>"},

		// Power off
		{"power-off", "udisksctl power-off -b /dev/sdb", "udisksctl power-off -b <path>"},

		// Info
		{"info block", "udisksctl info -b /dev/sda", "udisksctl info -b <path>"},

		// Fused flags
		{"mount fused block", "udisksctl mount --block-device=/dev/sda1", "udisksctl mount --block-device=<path>"},
		{"mount fused type", "udisksctl mount --filesystem-type=ext4 -b /dev/sda1", "udisksctl mount --filesystem-type=<val> -b <path>"},
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
		a := shellshape.Normalize("udisksctl mount -b /dev/sda1")
		b := shellshape.Normalize("udisksctl mount -b /dev/sdb2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("udisksctl mount -b /dev/sda1")
		subshell := shellshape.Normalize("udisksctl mount -b $(finddev)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
