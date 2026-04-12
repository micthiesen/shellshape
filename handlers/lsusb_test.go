package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLsusb(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare", "lsusb", "lsusb"},
		{"verbose", "lsusb -v", "lsusb -v"},
		{"tree", "lsusb -t", "lsusb -t"},
		{"bus device filter", "lsusb -s 001:003", "lsusb -s <val>"},
		{"vendor product filter", "lsusb -d 2dc8:3106", "lsusb -d <val>"},
		{"device file", "lsusb -D /dev/bus/usb/001/003", "lsusb -D <path>"},
		{"verbose with device", "lsusb -v -d 8087:0029", "lsusb -v -d <val>"},
		{"verbose with slot", "lsusb -v -s 002:001", "lsusb -v -s <val>"},
		{"version", "lsusb -V", "lsusb -V"},
		{"with redirect", "lsusb -t > /tmp/usb.txt", "lsusb -t > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("lsusb -d 2dc8:3106")
		b := shellshape.Normalize("lsusb -d 8087:0029")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lsusb -d 2dc8:3106")
		subshell := shellshape.Normalize("lsusb -d $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
