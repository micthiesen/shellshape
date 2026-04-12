package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestEfibootmgr(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "efibootmgr", "efibootmgr"},
		{"verbose", "efibootmgr -v", "efibootmgr -v"},

		// Delete boot entry
		{"delete entry", "efibootmgr -b 0003 -B", "efibootmgr -b N -B"},

		// Set boot order
		{"boot order", "efibootmgr -o 0001,0002,0003", "efibootmgr -o <val>"},

		// Set next boot
		{"next boot", "efibootmgr -n 0001", "efibootmgr -n N"},

		// Set timeout
		{"timeout", "efibootmgr -t 5", "efibootmgr -t N"},

		// Create entry
		{"create entry", "efibootmgr --create -d /dev/sda -p 1 -l /EFI/arch/grubx64.efi -L Arch", "efibootmgr --create -d <path> -p N -l <path> -L <val>"},

		// Combination
		{"create full", "efibootmgr -b 0004 --create -d /dev/nvme0n1 -p 1 -l /EFI/Linux/linux.efi -L Linux", "efibootmgr -b N --create -d <path> -p N -l <path> -L <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different boot entries collide", func(t *testing.T) {
		a := shellshape.Normalize("efibootmgr -b 0001 -B")
		b := shellshape.Normalize("efibootmgr -b 0004 -B")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("efibootmgr -b 0001")
		subshell := shellshape.Normalize("efibootmgr $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
