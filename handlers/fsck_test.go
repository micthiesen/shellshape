package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFsck(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single device", "fsck /dev/sda1", "fsck <path>"},
		{"no args", "fsck", "fsck"},

		// Boolean flags
		{"auto repair", "fsck -a /dev/sda1", "fsck -a <path>"},
		{"preen", "fsck -p /dev/sda1", "fsck -p <path>"},
		{"yes mode", "fsck -y /dev/sda1", "fsck -y <path>"},
		{"no mode", "fsck -n /dev/sda1", "fsck -n <path>"},
		{"force check", "fsck -f /dev/sda1", "fsck -f <path>"},
		{"verbose", "fsck -V /dev/sda1", "fsck -V <path>"},
		{"combined flags", "fsck -pf /dev/sda1", "fsck -pf <path>"},

		// -t flag consumes filesystem type
		{"type flag", "fsck -t ext4 /dev/sda1", "fsck -t <type> <path>"},
		{"type flag xfs", "fsck -t xfs /dev/vda1", "fsck -t <type> <path>"},
		{"type with other flags", "fsck -y -t ext4 /dev/sda1", "fsck -y -t <type> <path>"},

		// Multiple positionals
		{"multiple devices", "fsck /dev/sda1 /dev/sda2", "fsck <path>+"},
		{"multiple with flag", "fsck -a /dev/sda1 /dev/sda2 /dev/sdb1", "fsck -a <path>+"},

		// Redirects
		{"with redirect", "fsck /dev/sda1 > /tmp/fsck.log", "fsck <path> > <path>"},
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
	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("fsck /dev/sda1")
		b := shellshape.Normalize("fsck /dev/nvme0n1p3")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different filesystem types collapse to same shape
	t.Run("different types collide", func(t *testing.T) {
		a := shellshape.Normalize("fsck -t ext4 /dev/sda1")
		b := shellshape.Normalize("fsck -t xfs /dev/vda1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fsck /dev/sda1")
		subshell := shellshape.Normalize("fsck $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
