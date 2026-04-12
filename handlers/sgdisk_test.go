package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSgdisk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"zap all", "sgdisk -Z /dev/sda", "sgdisk -Z <path>"},
		{"print partitions", "sgdisk -p /dev/nvme0n1", "sgdisk -p <path>"},
		{"partition info", "sgdisk -i 1 /dev/sda", "sgdisk -i N <path>"},

		// Creating partitions
		{"new partition", "sgdisk -n 1:0:+512M /dev/sda", "sgdisk -n <val> <path>"},
		{"new with type", "sgdisk -n 1:0:+512M -t 1:ef00 /dev/sda", "sgdisk -n <val> -t <val> <path>"},
		{"new with name", "sgdisk -n 2:0:0 -t 2:8300 -c 2:root /dev/sda", "sgdisk -n <val> -t <val> -c <val> <path>"},

		// Delete partition
		{"delete partition", "sgdisk -d 3 /dev/sda", "sgdisk -d N <path>"},

		// Multiple operations
		{"zap and create", "sgdisk -Z -n 1:0:+512M -t 1:ef00 -n 2:0:0 -t 2:8300 /dev/sda", "sgdisk -Z -n <val> -t <val> -n <val> -t <val> <path>"},

		// Backup/restore
		{"backup fused", "sgdisk --backup=/tmp/sda.bak /dev/sda", "sgdisk --backup=<path> <path>"},
		{"load-backup fused", "sgdisk --load-backup=/tmp/sda.bak /dev/sda", "sgdisk --load-backup=<path> <path>"},

		// Other boolean flags
		{"sort partitions", "sgdisk -s /dev/sda", "sgdisk -s <path>"},
		{"randomize GUIDs", "sgdisk -G /dev/sdc", "sgdisk -G <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("sgdisk -Z /dev/sda")
		b := shellshape.Normalize("sgdisk -Z /dev/nvme0n1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different partition specs collide", func(t *testing.T) {
		a := shellshape.Normalize("sgdisk -n 1:0:+512M /dev/sda")
		b := shellshape.Normalize("sgdisk -n 2:0:+1G /dev/sda")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sgdisk literal-arg")
		subshell := shellshape.Normalize("sgdisk $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
