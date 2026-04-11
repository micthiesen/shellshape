package shellshape

import "testing"

func TestFdisk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list partitions", "fdisk -l", "fdisk -l"},
		{"list specific device", "fdisk -l /dev/sda", "fdisk -l <path>"},
		{"open device", "fdisk /dev/sda", "fdisk <path>"},
		{"partition size", "fdisk -s /dev/sda1", "fdisk -s <path>"},

		// Flags with numeric arguments
		{"sector size", "fdisk -b 4096 /dev/sda", "fdisk -b N <path>"},
		{"cylinders", "fdisk -C 1024 /dev/sda", "fdisk -C N <path>"},
		{"heads", "fdisk -H 255 /dev/sda", "fdisk -H N <path>"},
		{"sectors per track", "fdisk -S 63 /dev/sda", "fdisk -S N <path>"},
		{"multiple numeric flags", "fdisk -C 1024 -H 255 -S 63 /dev/sda", "fdisk -C N -H N -S N <path>"},

		// Boolean flags
		{"help flag", "fdisk -h", "fdisk -h"},
		{"version flag", "fdisk -v", "fdisk -v"},

		// Combined flags
		{"list with sector size", "fdisk -b 512 -l /dev/sdb", "fdisk -b N -l <path>"},

		// Flag arguments that generic classifier wouldn't handle correctly
		{"sector size non-numeric", "fdisk -b auto /dev/sda", "fdisk -b N <path>"},
		{"cylinders non-numeric", "fdisk -C max /dev/sda", "fdisk -C N <path>"},
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
	t.Run("different devices collide", func(t *testing.T) {
		a := Normalize("fdisk /dev/sda")
		b := Normalize("fdisk /dev/nvme0n1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric values collide", func(t *testing.T) {
		a := Normalize("fdisk -b 512 /dev/sda")
		b := Normalize("fdisk -b 4096 /dev/sdb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("fdisk /dev/sda")
		subshell := Normalize("fdisk $(echo /dev/sda)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
