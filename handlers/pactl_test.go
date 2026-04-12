package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPactl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"list sinks", "pactl list sinks", "pactl list sinks"},
		{"list sources short", "pactl list short sources", "pactl list short sources"},
		{"stat", "pactl stat", "pactl stat"},
		{"info", "pactl info", "pactl info"},

		// Volume control (sink + volume both <val>, collapsed to <val>+)
		{"set sink volume", "pactl set-sink-volume @DEFAULT_SINK@ 50%", "pactl set-sink-volume <val>+"},
		{"set source volume", "pactl set-source-volume alsa_input.usb +5%", "pactl set-source-volume <val>+"},
		{"set sink volume by index", "pactl set-sink-volume 0 75%", "pactl set-sink-volume <val>+"},

		// Mute control
		{"set sink mute", "pactl set-sink-mute @DEFAULT_SINK@ toggle", "pactl set-sink-mute <val>+"},
		{"set source mute", "pactl set-source-mute 1 0", "pactl set-source-mute <val>+"},

		// Move streams
		{"move sink input", "pactl move-sink-input 3 alsa_output.pci", "pactl move-sink-input <val>+"},

		// Module management
		{"load module", "pactl load-module module-null-sink sink_name=MySink", "pactl load-module module-null-sink <val>"},
		{"unload module", "pactl unload-module 42", "pactl unload-module <val>"},

		// Default sink/source
		{"set default sink", "pactl set-default-sink alsa_output.pci-0000_00_1f.3.analog-stereo", "pactl set-default-sink <val>"},
		{"set default source", "pactl set-default-source alsa_input.usb-Creative_Technology_Ltd", "pactl set-default-source <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test
	t.Run("different sink names collide", func(t *testing.T) {
		a := shellshape.Normalize("pactl set-sink-volume alsa_output.pci 50%")
		b := shellshape.Normalize("pactl set-sink-volume alsa_output.usb 75%")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pactl set-sink-volume @DEFAULT_SINK@ 50%")
		subshell := shellshape.Normalize("pactl set-sink-volume $(evil-cmd) 50%")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
