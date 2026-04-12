package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNvme(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"list", "nvme list", "nvme list"},
		{"smart-log", "nvme smart-log /dev/nvme0", "nvme smart-log <path>"},
		{"id-ctrl", "nvme id-ctrl /dev/nvme0", "nvme id-ctrl <path>"},
		{"id-ns", "nvme id-ns /dev/nvme0n1", "nvme id-ns <path>"},
		{"fw-log", "nvme fw-log /dev/nvme0", "nvme fw-log <path>"},
		{"error-log", "nvme error-log /dev/nvme0", "nvme error-log <path>"},

		// Namespace ID (numeric)
		{"id-ns with namespace", "nvme id-ns /dev/nvme0 -n 1", "nvme id-ns <path> -n N"},
		{"namespace long flag", "nvme id-ns /dev/nvme0 --namespace-id 2", "nvme id-ns <path> --namespace-id N"},

		// Output format (structural)
		{"output format", "nvme smart-log /dev/nvme0 -o json", "nvme smart-log <path> -o json"},
		{"output format long", "nvme smart-log /dev/nvme0 --output-format=json", "nvme smart-log <path> --output-format=json"},

		// Format command
		{"format", "nvme format /dev/nvme0n1 -l 0", "nvme format <path> -l N"},
		{"format long", "nvme format /dev/nvme0n1 --lba-format 1", "nvme format <path> --lba-format N"},
		{"format with ses", "nvme format /dev/nvme0 -s 1 -n 1", "nvme format <path> -s N -n N"},

		// Boolean flags
		{"human readable", "nvme smart-log /dev/nvme0 -H", "nvme smart-log <path> -H"},
		{"verbose", "nvme smart-log /dev/nvme0 -v", "nvme smart-log <path> -v"},

		// Get/set feature
		{"get-feature", "nvme get-feature /dev/nvme0 -f 2", "nvme get-feature <path> -f N"},
		{"set-feature", "nvme set-feature /dev/nvme0 -f 2 -V 1", "nvme set-feature <path> -f N -V N"},

		// Sanitize
		{"sanitize", "nvme sanitize /dev/nvme0 -a 2", "nvme sanitize <path> -a N"},

		// Get-log
		{"get-log", "nvme get-log /dev/nvme0 -i 2 -l 512", "nvme get-log <path> -i N -l N"},
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
		a := shellshape.Normalize("nvme smart-log /dev/nvme0")
		b := shellshape.Normalize("nvme smart-log /dev/nvme1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nvme smart-log /dev/nvme0")
		subshell := shellshape.Normalize("nvme smart-log $(echo /dev/nvme0)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
