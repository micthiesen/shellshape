package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNvidiaSmi(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "nvidia-smi", "nvidia-smi"},
		{"query flag", "nvidia-smi -q", "nvidia-smi -q"},
		{"list gpus", "nvidia-smi -L", "nvidia-smi -L"},

		// Display type (structural, kept verbatim)
		{"display memory", "nvidia-smi -q -d MEMORY", "nvidia-smi -q -d MEMORY"},
		{"display temperature", "nvidia-smi -q -d TEMPERATURE", "nvidia-smi -q -d TEMPERATURE"},

		// Numeric flags
		{"gpu id", "nvidia-smi -i 0 -q", "nvidia-smi -i N -q"},
		{"loop interval", "nvidia-smi -l 5", "nvidia-smi -l N"},
		{"fused id", "nvidia-smi --id=2", "nvidia-smi --id=N"},

		// Query GPU (collapse value)
		{"query gpu fields", "nvidia-smi --query-gpu=index,name,uuid,serial --format=csv", "nvidia-smi --query-gpu=<val> --format=csv"},
		{"query gpu other fields", "nvidia-smi --query-gpu=temperature.gpu,power.draw --format=csv,noheader,nounits", "nvidia-smi --query-gpu=<val> --format=csv,noheader,nounits"},

		// nvidia-settings alias
		{"nvidia-settings", "nvidia-settings --assign", "nvidia-settings --assign"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different query-gpu fields collide", func(t *testing.T) {
		a := shellshape.Normalize("nvidia-smi --query-gpu=index,name --format=csv")
		b := shellshape.Normalize("nvidia-smi --query-gpu=temperature.gpu,power.draw --format=csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nvidia-smi -i 0")
		subshell := shellshape.Normalize("nvidia-smi $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
