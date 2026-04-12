package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSensors(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare", "sensors", "sensors"},
		{"fahrenheit", "sensors -f", "sensors -f"},
		{"json", "sensors -j", "sensors -j"},
		{"raw", "sensors -u", "sensors -u"},
		{"no adapter", "sensors -A", "sensors -A"},
		{"chip name", "sensors coretemp-isa-0000", "sensors <val>"},
		{"chip with flag", "sensors -u coretemp-isa-0000", "sensors -u <val>"},
		{"config file", "sensors -c /etc/sensors3.conf", "sensors -c <path>"},
		{"config file long", "sensors --config-file /etc/sensors.conf", "sensors --config-file <path>"},
		{"bus filter", "sensors --bus 0", "sensors --bus <val>"},
		{"multiple chips", "sensors coretemp-isa-0000 k10temp-pci-00c3", "sensors <val>+"},
		{"with redirect", "sensors -j > /tmp/sensors.json", "sensors -j > <path>"},
		{"sensors-detect bare", "sensors-detect", "sensors-detect"},
		{"sensors-detect auto", "sensors-detect --auto", "sensors-detect --auto"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different chips collide", func(t *testing.T) {
		a := shellshape.Normalize("sensors coretemp-isa-0000")
		b := shellshape.Normalize("sensors k10temp-pci-00c3")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("sensors coretemp-isa-0000")
		subshell := shellshape.Normalize("sensors $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
