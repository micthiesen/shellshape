package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPicocom(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple device", "picocom /dev/ttyUSB0", "picocom <path>"},
		{"device with baud", "picocom -b 115200 /dev/ttyUSB0", "picocom -b N <path>"},
		{"long baud", "picocom --baud 9600 /dev/ttyACM0", "picocom --baud N <path>"},

		// Flow control and parity (structural)
		{"flow control", "picocom --flow hard /dev/ttyUSB0", "picocom --flow hard <path>"},
		{"parity", "picocom --parity even /dev/ttyS0", "picocom --parity even <path>"},

		// Numeric flags
		{"databits", "picocom --databits 8 /dev/ttyUSB0", "picocom --databits N <path>"},
		{"stopbits", "picocom --stopbits 1 /dev/ttyUSB0", "picocom --stopbits N <path>"},

		// Map flags
		{"imap", "picocom --imap crcrlf /dev/ttyUSB0", "picocom --imap <val> <path>"},
		{"omap", "picocom --omap delbs,crlf /dev/ttyUSB0", "picocom --omap <val> <path>"},
		{"emap", "picocom --emap crcrlf,delbs /dev/ttyUSB0", "picocom --emap <val> <path>"},

		// Boolean flags
		{"noreset", "picocom --noreset /dev/ttyUSB0", "picocom --noreset <path>"},
		{"noinit", "picocom --noinit /dev/ttyUSB0", "picocom --noinit <path>"},
		{"quiet", "picocom -q /dev/ttyUSB0", "picocom -q <path>"},

		// Combined
		{"full config", "picocom -b 115200 --flow none --parity none --databits 8 --stopbits 1 --noreset /dev/ttyUSB0", "picocom -b N --flow none --parity none --databits N --stopbits N --noreset <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("picocom /dev/ttyUSB0")
		b := shellshape.Normalize("picocom /dev/ttyACM1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different baud rates collide", func(t *testing.T) {
		a := shellshape.Normalize("picocom -b 9600 /dev/ttyUSB0")
		b := shellshape.Normalize("picocom -b 115200 /dev/ttyUSB0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("picocom literal-arg")
		subshell := shellshape.Normalize("picocom $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
