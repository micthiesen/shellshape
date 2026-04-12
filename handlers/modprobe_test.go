package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestModprobe(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - module names are structural
		{"load module", "modprobe nvidia", "modprobe nvidia"},
		{"remove module", "modprobe -r nvidia", "modprobe -r nvidia"},
		{"remove long", "modprobe --remove snd_hda_intel", "modprobe --remove snd_hda_intel"},
		{"dry-run", "modprobe -n btusb", "modprobe -n btusb"},
		{"verbose load", "modprobe -v iwlwifi", "modprobe -v iwlwifi"},
		{"quiet", "modprobe -q nouveau", "modprobe -q nouveau"},
		{"first-time", "modprobe --first-time vfio-pci", "modprobe --first-time vfio-pci"},
		{"show-depends", "modprobe --show-depends nvidia", "modprobe --show-depends nvidia"},

		// Module parameters collapse to <val>
		{"module with param", "modprobe nvidia NVreg_PreserveVideoMemoryAllocations=1", "modprobe nvidia <val>"},
		{"module with multiple params", "modprobe snd_hda_intel power_save=1 power_save_controller=Y", "modprobe snd_hda_intel <val>+"},
		{"module with complex param", "modprobe vfio-pci ids=10de:2684,10de:22bc", "modprobe vfio-pci <val>"},

		// Path flags
		{"config flag", "modprobe -C /etc/modprobe.d/custom.conf nvidia", "modprobe -C <path> nvidia"},
		{"directory flag", "modprobe -d /lib/modules nvidia", "modprobe -d <path> nvidia"},

		// Combined flags
		{"verbose remove", "modprobe -rv nvidia", "modprobe -rv nvidia"},
		{"multiple flags", "modprobe -v -n btusb", "modprobe -v -n btusb"},
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
	t.Run("different params collide", func(t *testing.T) {
		a := shellshape.Normalize("modprobe nvidia NVreg_PreserveVideoMemoryAllocations=1")
		b := shellshape.Normalize("modprobe nvidia NVreg_PreserveVideoMemoryAllocations=0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("modprobe nvidia param=value")
		subshell := shellshape.Normalize("modprobe nvidia $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
