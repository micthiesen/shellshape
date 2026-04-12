package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBtrfs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Subvolume sub-subcommands
		{"subvolume create", "btrfs subvolume create /mnt/sv1", "btrfs subvolume create <path>"},
		{"subvolume list", "btrfs subvolume list /mnt", "btrfs subvolume list <path>"},
		{"subvolume delete", "btrfs subvolume delete /mnt/sv1", "btrfs subvolume delete <path>"},
		{"subvolume snapshot", "btrfs subvolume snapshot /mnt/sv1 /mnt/snap", "btrfs subvolume snapshot <path>+"},
		{"subvolume show", "btrfs subvolume show /mnt/sv1", "btrfs subvolume show <path>"},
		{"subvolume set-default", "btrfs subvolume set-default 256 /mnt", "btrfs subvolume set-default N <path>"},

		// Filesystem sub-subcommands
		{"filesystem df", "btrfs filesystem df /mnt", "btrfs filesystem df <path>"},
		{"filesystem show", "btrfs filesystem show /dev/sda", "btrfs filesystem show <path>"},
		{"filesystem resize", "btrfs filesystem resize max /mnt", "btrfs filesystem resize max <path>"},
		{"filesystem usage", "btrfs filesystem usage /mnt", "btrfs filesystem usage <path>"},

		// Device sub-subcommands
		{"device add", "btrfs device add /dev/sdf /mnt", "btrfs device add <path>+"},
		{"device delete", "btrfs device delete /dev/sdf /mnt", "btrfs device delete <path>+"},
		{"device stats", "btrfs device stats /mnt", "btrfs device stats <path>"},

		// Balance
		{"balance start", "btrfs balance start /mnt", "btrfs balance start <path>"},
		{"balance status", "btrfs balance status /mnt", "btrfs balance status <path>"},

		// Scrub
		{"scrub start", "btrfs scrub start /mnt", "btrfs scrub start <path>"},
		{"scrub status", "btrfs scrub status /mnt", "btrfs scrub status <path>"},

		// Send/receive
		{"send", "btrfs send /mnt/snap", "btrfs send <path>"},
		{"receive", "btrfs receive /mnt/dest", "btrfs receive <path>"},
		{"send with parent", "btrfs send -p /mnt/parent /mnt/snap", "btrfs send -p <path>+"},

		// Property
		{"property set", "btrfs property set /mnt/sv1 ro true", "btrfs property set <path> ro true"},
		{"property get", "btrfs property get /mnt/sv1", "btrfs property get <path>"},

		// Quota/qgroup
		{"quota enable", "btrfs quota enable /mnt", "btrfs quota enable <path>"},
		{"qgroup show", "btrfs qgroup show /mnt", "btrfs qgroup show <path>"},

		// Check
		{"check", "btrfs check /dev/sda1", "btrfs check <path>"},
		{"check readonly", "btrfs check --readonly /dev/sda1", "btrfs check --readonly <path>"},

		// Flags with values
		{"subvolume list with flags", "btrfs subvolume list -t /mnt", "btrfs subvolume list -t <path>"},
		{"subvolume delete with commit-after", "btrfs subvolume delete --commit-after /mnt/sv1", "btrfs subvolume delete --commit-after <path>"},
		{"inspect-internal rootid", "btrfs inspect-internal rootid /mnt", "btrfs inspect-internal rootid <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("btrfs subvolume create /mnt/vol1")
		b := shellshape.Normalize("btrfs subvolume create /mnt/vol2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("btrfs device add /dev/sdb /mnt")
		b := shellshape.Normalize("btrfs device add /dev/sdc /data")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("btrfs subvolume create /mnt/sv1")
		subshell := shellshape.Normalize("btrfs subvolume create $(echo /mnt/sv1)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
