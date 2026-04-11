package shellshape

import "testing"

func TestParted(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list all", "parted --list", "parted --list"},
		{"list short", "parted -l", "parted -l"},
		{"print partitions", "parted /dev/sda print", "parted <device> print"},
		{"print free space", "parted /dev/sda print free", "parted <device> print free"},

		// Flags with arguments
		{"align flag", "parted -a optimal /dev/sda print", "parted -a <val> <device> print"},
		{"align long flag", "parted --align optimal /dev/sda print", "parted --align <val> <device> print"},
		{"script mode", "parted --script /dev/sda mklabel gpt", "parted --script <device> mklabel <val>"},

		// Inline commands with arguments
		{"mklabel", "parted /dev/sda mklabel gpt", "parted <device> mklabel <val>"},
		{"mkpart basic", "parted /dev/sda mkpart primary 0% 100%", "parted <device> mkpart <val>+"},
		{"mkpart with fs", "parted /dev/sda mkpart primary ext4 0% 50%", "parted <device> mkpart <val>+"},
		{"rm partition", "parted /dev/sda rm 1", "parted <device> rm <val>"},
		{"resizepart", "parted /dev/sda resizepart 1 100%", "parted <device> resizepart <val>+"},
		{"set flag", "parted /dev/sda set 1 boot on", "parted <device> set <val>+"},
		{"toggle", "parted /dev/sda toggle 1 boot", "parted <device> toggle <val>+"},
		{"name partition", "parted /dev/sda name 1 mypart", "parted <device> name <val>+"},
		{"unit then print", "parted /dev/sda unit GB print", "parted <device> unit <val> print"},
		{"select device", "parted /dev/sda select /dev/sdb", "parted <device> select <val>"},
		{"rescue", "parted /dev/sda rescue 0 100%", "parted <device> rescue <val>+"},
		{"align-check", "parted /dev/sda align-check optimal 1", "parted <device> align-check <val>+"},

		// Multiple flags combined
		{"script and align", "parted -s -a optimal /dev/sda mkpart primary 0% 100%", "parted -s -a <val> <device> mkpart <val>+"},

		// Device path variations
		{"nvme device", "parted /dev/nvme0n1 print", "parted <device> print"},
		{"disk by-id", "parted /dev/disk/by-id/scsi-0 print", "parted <device> print"},
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
		a := Normalize("parted /dev/sda print")
		b := Normalize("parted /dev/nvme0n1 print")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different label types collide", func(t *testing.T) {
		a := Normalize("parted /dev/sda mklabel gpt")
		b := Normalize("parted /dev/sda mklabel msdos")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("parted /dev/sda print")
		subshell := Normalize("parted $(dangerous-command) print")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
