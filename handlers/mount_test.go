package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "mount", "mount"},
		{"single path", "mount /mnt", "mount <path>"},
		{"device and target", "mount /dev/sda1 /mnt", "mount <path>+"},

		// Flags with arguments
		{"type flag short", "mount -t ext4 /dev/sda1 /mnt", "mount -t <type> <path>+"},
		{"type flag long", "mount --types nfs server:/share /data", "mount --types <type> <path>+"},
		{"options flag", "mount -o rw,noexec /dev/sdb1 /data", "mount -o <opts> <path>+"},
		{"options flag long", "mount --options rw,nosuid /dev/sdc /srv", "mount --options <opts> <path>+"},
		{"label flag", "mount -L BACKUP /backup", "mount -L <label> <path>"},
		{"label flag long", "mount --label MYDATA /data", "mount --label <label> <path>"},
		{"uuid flag", "mount -U 550e8400-e29b /mnt", "mount -U <uuid> <path>"},
		{"uuid flag long", "mount --uuid abc123 /mnt", "mount --uuid <uuid> <path>"},
		{"fstab flag", "mount -T /etc/fstab.alt -a", "mount -T <path> -a"},
		{"fstab flag long", "mount --fstab /etc/fstab.alt -a", "mount --fstab <path> -a"},
		{"namespace flag", "mount -N 1234 /dev/sda1 /mnt", "mount -N N <path>+"},
		{"source flag", "mount --source /dev/sda1 --target /mnt", "mount --source <path> --target <path>"},
		{"test-opts flag", "mount -O no_netdev -a", "mount -O <opts> -a"},

		// Boolean flags
		{"all", "mount -a", "mount -a"},
		{"read-only", "mount -r /dev/cdrom /cdrom", "mount -r <path>+"},
		{"bind", "mount --bind /src /dst", "mount --bind <path>+"},
		{"verbose", "mount -v -t ext4 /dev/sda1 /mnt", "mount -v -t <type> <path>+"},

		// Combined flags and options
		{"type and options", "mount -t tmpfs -o size=4G tmpfs /mnt", "mount -t <type> -o <opts> tmpfs <path>"},
		{"ro with type", "mount -t iso9660 -o ro /dev/cdrom /cdrom", "mount -t <type> -o <opts> <path>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different device paths collapse to same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("mount /dev/sda1 /mnt/disk1")
		b := shellshape.Normalize("mount /dev/sdb2 /mnt/disk2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different mount options collapse to same shape
	t.Run("different options collide", func(t *testing.T) {
		a := shellshape.Normalize("mount -o rw,noexec /dev/sda1 /data")
		b := shellshape.Normalize("mount -o ro,nosuid /dev/sdb1 /srv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mount /dev/sda1 /mnt")
		subshell := shellshape.Normalize("mount $(dangerous-command) /mnt")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestUmount(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single path", "umount /mnt", "umount <path>"},
		{"device path", "umount /dev/sda1", "umount <path>"},

		// Flags with arguments
		{"type flag", "umount -t nfs -a", "umount -t <type> -a"},
		{"namespace flag", "umount -N 5678", "umount -N N"},
		{"test-opts flag", "umount -O no_netdev -a", "umount -O <opts> -a"},

		// Boolean flags
		{"all", "umount -a", "umount -a"},
		{"force", "umount -f /mnt/nfs", "umount -f <path>"},
		{"lazy", "umount -l /mnt", "umount -l <path>"},
		{"recursive", "umount -R /mnt", "umount -R <path>"},

		// Multiple paths
		{"multiple paths", "umount /mnt/a /mnt/b /mnt/c", "umount <path>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different mount points collapse
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("umount /mnt/disk1")
		b := shellshape.Normalize("umount /mnt/disk2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("umount /mnt")
		subshell := shellshape.Normalize("umount $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
