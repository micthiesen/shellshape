package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestQuickemu(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic vm", "quickemu --vm /vms/windows-11.conf", "quickemu --vm <path>"},
		{"display spice", "quickemu --vm /vms/ubuntu.conf --display spice", "quickemu --vm <path> --display spice"},
		{"fullscreen", "quickemu --fullscreen --vm /vms/macos.conf", "quickemu --fullscreen --vm <path>"},
		{"snapshot create", "quickemu --snapshot create mytag --vm /vms/windows-11.conf", "quickemu --snapshot create <val> --vm <path>"},
		{"snapshot apply", "quickemu --snapshot apply backup1 --vm /vms/windows-11.conf", "quickemu --snapshot apply <val> --vm <path>"},
		{"snapshot delete", "quickemu --snapshot delete old --vm /vms/windows-11.conf", "quickemu --snapshot delete <val> --vm <path>"},
		{"snapshot info", "quickemu --snapshot info --vm /vms/windows-11.conf", "quickemu --snapshot info --vm <path>"},
		{"ssh port", "quickemu --vm /vms/ubuntu.conf --ssh-port 2222", "quickemu --vm <path> --ssh-port <val>"},
		{"status quo", "quickemu --status-quo --vm /vms/windows-11.conf", "quickemu --status-quo --vm <path>"},
		{"public dir", "quickemu --vm /vms/ubuntu.conf --public-dir /home/user/Public", "quickemu --vm <path> --public-dir <path>"},
		{"keyboard and mouse", "quickemu --vm /vms/ubuntu.conf --keyboard virtio --mouse tablet", "quickemu --vm <path> --keyboard virtio --mouse tablet"},
		{"multiple options", "quickemu --fullscreen --display gtk --sound-card intel-hda --vm /vms/windows-11.conf", "quickemu --fullscreen --display gtk --sound-card intel-hda --vm <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different VMs collide", func(t *testing.T) {
		a := shellshape.Normalize("quickemu --vm /vms/windows-11.conf")
		b := shellshape.Normalize("quickemu --vm /vms/ubuntu-22.04.conf")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("quickemu --vm /vms/windows-11.conf")
		subshell := shellshape.Normalize("quickemu --vm $(echo evil)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
