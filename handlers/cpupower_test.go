package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestCpupower(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"frequency-info", "cpupower frequency-info", "cpupower frequency-info"},
		{"frequency-set", "cpupower frequency-set", "cpupower frequency-set"},
		{"idle-info", "cpupower idle-info", "cpupower idle-info"},
		{"idle-set", "cpupower idle-set", "cpupower idle-set"},
		{"monitor", "cpupower monitor", "cpupower monitor"},
		{"info", "cpupower info", "cpupower info"},

		// Governor (structural)
		{"set governor", "cpupower frequency-set -g powersave", "cpupower frequency-set -g powersave"},
		{"set governor long", "cpupower frequency-set --governor performance", "cpupower frequency-set --governor performance"},

		// Frequency values (data)
		{"set freq", "cpupower frequency-set -f 2.4GHz", "cpupower frequency-set -f <val>"},
		{"set freq long", "cpupower frequency-set --freq 3000MHz", "cpupower frequency-set --freq <val>"},
		{"set min", "cpupower frequency-set -d 1GHz", "cpupower frequency-set -d <val>"},
		{"set max", "cpupower frequency-set -u 3.5GHz", "cpupower frequency-set -u <val>"},

		// --cpu flag before subcommand
		{"cpu all freq-info", "cpupower --cpu all frequency-info", "cpupower --cpu <val> frequency-info"},
		{"cpu 0 freq-set", "cpupower --cpu 0 frequency-set -g performance", "cpupower --cpu <val> frequency-set -g performance"},
		{"cpu range", "cpupower --cpu 0-3 frequency-set -g schedutil", "cpupower --cpu <val> frequency-set -g schedutil"},
		{"-c shorthand", "cpupower -c all info", "cpupower -c <val> info"},

		// Boolean flags
		{"human", "cpupower frequency-info --human", "cpupower frequency-info --human"},
		{"hwfreq", "cpupower frequency-info --hwfreq", "cpupower frequency-info --hwfreq"},

		// Idle set
		{"idle disable", "cpupower idle-set -d 3", "cpupower idle-set -d <val>"},
		{"idle enable", "cpupower idle-set -e 2", "cpupower idle-set -e <val>"},
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
	t.Run("different frequencies collide", func(t *testing.T) {
		a := shellshape.Normalize("cpupower frequency-set -f 2.4GHz")
		b := shellshape.Normalize("cpupower frequency-set -f 3.0GHz")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different cpu selectors collide", func(t *testing.T) {
		a := shellshape.Normalize("cpupower --cpu 0 frequency-info")
		b := shellshape.Normalize("cpupower --cpu all frequency-info")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("cpupower frequency-set -f 2GHz")
		subshell := shellshape.Normalize("cpupower frequency-set -f $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
