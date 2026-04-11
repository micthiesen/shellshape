package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBlkid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "blkid", "blkid"},
		{"single device", "blkid /dev/sda1", "blkid <path>"},
		{"multiple devices", "blkid /dev/sda1 /dev/sdb1", "blkid <path>+"},

		// Flags with value arguments
		{"tag match", "blkid -s UUID -o value /dev/sda1", "blkid -s <val> -o <val> <path>"},
		{"label lookup", "blkid -L mylabel", "blkid -L <val>"},
		{"uuid lookup", "blkid -U 1234-5678", "blkid -U <val>"},
		{"token match", "blkid -t TYPE=ext4", "blkid -t <val>"},
		{"type restriction", "blkid -n ext4,xfs /dev/sda1", "blkid -n <val> <path>"},
		{"usage restriction", "blkid -u filesystem,raid /dev/sda", "blkid -u <val> <path>"},

		// Cache file (path argument)
		{"cache file", "blkid -c /dev/null /dev/sda1", "blkid -c <path>+"},

		// Numeric arguments
		{"offset", "blkid -O 1024 /dev/sda", "blkid -O N <path>"},
		{"size", "blkid -S 512 /dev/sda", "blkid -S N <path>"},

		// Boolean flags
		{"probe mode", "blkid -p -o udev /dev/sda1", "blkid -p -o <val> <path>"},
		{"garbage collect", "blkid -g", "blkid -g"},
		{"list one", "blkid -l -t TYPE=ext4", "blkid -l -t <val>"},
		{"info flag", "blkid -p -i /dev/sda", "blkid -p -i <path>"},

		// Long flags
		{"long output", "blkid --output value /dev/sda1", "blkid --output <val> <path>"},
		{"long cache", "blkid --cache-file /tmp/cache /dev/sda", "blkid --cache-file <path>+"},
		{"long match-tag", "blkid --match-tag UUID /dev/sda", "blkid --match-tag <val> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different device paths produce the same shape
	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("blkid -s UUID /dev/sda1")
		b := shellshape.Normalize("blkid -s UUID /dev/nvme0n1p2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different labels produce the same shape
	t.Run("different labels collide", func(t *testing.T) {
		a := shellshape.Normalize("blkid -L boot")
		b := shellshape.Normalize("blkid -L rootfs")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("blkid /dev/sda1")
		subshell := shellshape.Normalize("blkid $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
