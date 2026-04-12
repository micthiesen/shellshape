package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBootctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"status", "bootctl status", "bootctl status"},
		{"install", "bootctl install", "bootctl install"},
		{"update", "bootctl update", "bootctl update"},
		{"remove", "bootctl remove", "bootctl remove"},
		{"list", "bootctl list", "bootctl list"},
		{"is-installed", "bootctl is-installed", "bootctl is-installed"},
		{"random-seed", "bootctl random-seed", "bootctl random-seed"},

		// Subcommands with positional arguments
		{"set-default", "bootctl set-default auto-windows", "bootctl set-default <val>"},
		{"set-oneshot", "bootctl set-oneshot auto-reboot-to-firmware-setup", "bootctl set-oneshot <val>"},
		{"set-timeout", "bootctl set-timeout 5", "bootctl set-timeout N"},
		{"set-timeout menu", "bootctl set-timeout menu-force", "bootctl set-timeout <val>"},

		// Path flags after subcommand
		{"esp-path after sub", "bootctl install --esp-path /efi", "bootctl install --esp-path <path>"},
		{"esp-path fused after sub", "bootctl install --esp-path=/efi", "bootctl install --esp-path=<path>"},
		{"boot-path after sub", "bootctl status --boot-path /boot", "bootctl status --boot-path <path>"},
		{"boot-path fused after sub", "bootctl status --boot-path=/boot", "bootctl status --boot-path=<path>"},

		// Path flags before subcommand (separate form)
		{"esp-path before sub", "bootctl --esp-path /efi status", "bootctl --esp-path <path> status"},
		{"boot-path before sub", "bootctl --boot-path /boot list", "bootctl --boot-path <path> list"},

		// Boolean flags
		{"no-pager", "bootctl list --no-pager", "bootctl list --no-pager"},
		{"no-variables", "bootctl install --no-variables", "bootctl install --no-variables"},
		{"graceful", "bootctl update --graceful", "bootctl update --graceful"},

		// Entry token flag
		{"entry-token", "bootctl install --entry-token machine-id", "bootctl install --entry-token <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different entry IDs collide", func(t *testing.T) {
		a := shellshape.Normalize("bootctl set-default auto-windows")
		b := shellshape.Normalize("bootctl set-default linux-mainline.conf")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different esp paths collide", func(t *testing.T) {
		a := shellshape.Normalize("bootctl install --esp-path /efi")
		b := shellshape.Normalize("bootctl install --esp-path /boot/efi")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("bootctl set-default myentry")
		subshell := shellshape.Normalize("bootctl set-default $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
