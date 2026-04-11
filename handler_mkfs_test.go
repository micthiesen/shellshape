package shellshape

import "testing"

func TestMkfs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare mkfs", "mkfs /dev/sda1", "mkfs <path>"},
		{"mkfs with type flag", "mkfs -t ext4 /dev/sda1", "mkfs -t ext4 <path>"},
		{"mkfs with type and check", "mkfs -c -t ntfs /dev/sda1", "mkfs -c -t ntfs <path>"},

		// mkfs.ext4 variant
		{"ext4 basic", "mkfs.ext4 /dev/sda1", "mkfs.ext4 <path>"},
		{"ext4 with label", "mkfs.ext4 -L my-data /dev/sda1", "mkfs.ext4 -L <label> <path>"},
		{"ext4 block size", "mkfs.ext4 -b 4096 /dev/sda1", "mkfs.ext4 -b N <path>"},
		{"ext4 inode size", "mkfs.ext4 -I 256 /dev/sda1", "mkfs.ext4 -I N <path>"},
		{"ext4 bytes per inode", "mkfs.ext4 -i 8192 /dev/sda1", "mkfs.ext4 -i N <path>"},
		{"ext4 num inodes", "mkfs.ext4 -N 100000 /dev/sda1", "mkfs.ext4 -N N <path>"},
		{"ext4 reserved pct", "mkfs.ext4 -m 1 /dev/sda1", "mkfs.ext4 -m N <path>"},
		{"ext4 uuid", "mkfs.ext4 -U 550e8400-e29b-41d4-a716-446655440000 /dev/sda1", "mkfs.ext4 -U <uuid> <path>"},
		{"ext4 features", "mkfs.ext4 -O has_journal,extent /dev/sda1", "mkfs.ext4 -O has_journal,extent <path>"},
		{"ext4 force", "mkfs.ext4 -F /dev/sda1", "mkfs.ext4 -F <path>"},
		{"ext4 combined", "mkfs.ext4 -b 4096 -i 8192 -L backup /dev/sdc1", "mkfs.ext4 -b N -i N -L <label> <path>"},

		// mkfs.xfs variant
		{"xfs basic", "mkfs.xfs /dev/sdb1", "mkfs.xfs <path>"},
		{"xfs force", "mkfs.xfs -f /dev/sdb1", "mkfs.xfs -f <path>"},
		{"xfs block size", "mkfs.xfs -b size=4096 /dev/sdb1", "mkfs.xfs -b size=N <path>"},
		{"xfs label", "mkfs.xfs -L myvolume /dev/sdb1", "mkfs.xfs -L <label> <path>"},

		// mkfs.btrfs variant
		{"btrfs basic", "mkfs.btrfs /dev/sda1", "mkfs.btrfs <path>"},
		{"btrfs label", "mkfs.btrfs -L data /dev/sda1", "mkfs.btrfs -L <label> <path>"},
		{"btrfs force", "mkfs.btrfs -f /dev/sda1", "mkfs.btrfs -f <path>"},
		{"btrfs multiple devices", "mkfs.btrfs -L data /dev/sda /dev/sdb", "mkfs.btrfs -L <label> <path>+"},

		// mkfs.vfat variant
		{"vfat basic", "mkfs.vfat /dev/sda1", "mkfs.vfat <path>"},
		{"vfat fat32", "mkfs.vfat -F 32 /dev/sda1", "mkfs.vfat -F N <path>"},
		{"vfat label", "mkfs.vfat -n BOOT /dev/sda1", "mkfs.vfat -n <label> <path>"},

		// Long flags
		{"long type", "mkfs --type ext4 /dev/sda1", "mkfs --type ext4 <path>"},

		// Boolean flags preserved
		{"verbose", "mkfs.ext4 -v /dev/sda1", "mkfs.ext4 -v <path>"},
		{"quiet check", "mkfs.ext4 -q -c /dev/sda1", "mkfs.ext4 -q -c <path>"},

		// Redirect
		{"with redirect", "mkfs.ext4 /dev/sda1 2>&1", "mkfs.ext4 <path> 2>&1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different device paths collapse to same shape
	t.Run("different devices collide", func(t *testing.T) {
		a := Normalize("mkfs.ext4 /dev/sda1")
		b := Normalize("mkfs.ext4 /dev/nvme0n1p2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different labels collapse to same shape
	t.Run("different labels collide", func(t *testing.T) {
		a := Normalize("mkfs.ext4 -L my-data /dev/sda1")
		b := Normalize("mkfs.ext4 -L backups /dev/sda1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different numeric values collapse
	t.Run("different block sizes collide", func(t *testing.T) {
		a := Normalize("mkfs.ext4 -b 1024 /dev/sda1")
		b := Normalize("mkfs.ext4 -b 4096 /dev/sda1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("mkfs.ext4 /dev/sda1")
		subshell := Normalize("mkfs.ext4 $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
