package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMkinitcpio(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "mkinitcpio", "mkinitcpio"},
		{"allpresets", "mkinitcpio --allpresets", "mkinitcpio --allpresets"},
		{"listhooks", "mkinitcpio -L", "mkinitcpio -L"},

		// Preset (structural)
		{"preset short", "mkinitcpio -p linux", "mkinitcpio -p linux"},
		{"preset long", "mkinitcpio --preset linux-lts", "mkinitcpio --preset linux-lts"},

		// Path flags
		{"generate short", "mkinitcpio -g /boot/initramfs-linux.img", "mkinitcpio -g <path>"},
		{"generate long", "mkinitcpio --generate /boot/initramfs-linux.img", "mkinitcpio --generate <path>"},
		{"config short", "mkinitcpio -c /etc/mkinitcpio.conf", "mkinitcpio -c <path>"},
		{"config long", "mkinitcpio --config /etc/mkinitcpio.conf", "mkinitcpio --config <path>"},
		{"builddir", "mkinitcpio -t /tmp/mkinitcpio", "mkinitcpio -t <path>"},

		// Kernel version (data)
		{"kernel short", "mkinitcpio -k 6.1.0-arch1-1", "mkinitcpio -k <val>"},
		{"kernel long", "mkinitcpio --kernel 6.1.0-arch1-1", "mkinitcpio --kernel <val>"},

		// Hooks (structural)
		{"addhooks short", "mkinitcpio -A encrypt,lvm2", "mkinitcpio -A encrypt,lvm2"},
		{"addhooks long", "mkinitcpio --addhooks btrfs", "mkinitcpio --addhooks btrfs"},
		{"skiphooks", "mkinitcpio -S autodetect", "mkinitcpio -S autodetect"},
		{"hookhelp", "mkinitcpio -H encrypt", "mkinitcpio -H encrypt"},

		// Complex invocations
		{"custom generate", "mkinitcpio -c /etc/mkinitcpio.conf -g /boot/initramfs-linux.img -k 6.1.0",
			"mkinitcpio -c <path> -g <path> -k <val>"},
		{"preset with verbose", "mkinitcpio -p linux -v", "mkinitcpio -p linux -v"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different kernel versions collide", func(t *testing.T) {
		a := shellshape.Normalize("mkinitcpio -k 6.1.0-arch1-1")
		b := shellshape.Normalize("mkinitcpio -k 6.6.3-zen1-1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output paths collide", func(t *testing.T) {
		a := shellshape.Normalize("mkinitcpio -g /boot/initramfs-linux.img")
		b := shellshape.Normalize("mkinitcpio -g /boot/initramfs-linux-lts.img")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mkinitcpio -k 6.1.0")
		subshell := shellshape.Normalize("mkinitcpio -k $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
