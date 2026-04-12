package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLms(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Server subcommands
		{"server start", "lms server start", "lms server start"},
		{"server stop", "lms server stop", "lms server stop"},
		{"server status", "lms server status", "lms server status"},
		{"server start with port", "lms server start --port 8080", "lms server start --port N"},

		// Model operations
		{"load model", "lms load mlx-community/Llama-3-8B-Instruct", "lms load <path>"},
		{"unload model", "lms unload mlx-community/Llama-3-8B-Instruct", "lms unload <path>"},
		{"get model", "lms get mlx-community/Llama-3-8B-Instruct", "lms get <path>"},

		// List and status
		{"ls", "lms ls", "lms ls"},
		{"ps", "lms ps", "lms ps"},
		{"version", "lms version", "lms version"},
		{"log", "lms log", "lms log"},

		// GPU flag
		{"load with gpu", "lms load some-model --gpu max", "lms load <val> --gpu <val>"},
		{"load with gpu auto", "lms load some-model --gpu auto", "lms load <val> --gpu <val>"},

		// Create
		{"create with path", "lms create --gguf /path/to/model.gguf", "lms create --gguf <path>"},

		// Context length
		{"load with context", "lms load some-model --context-length 4096", "lms load <val> --context-length N"},
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
	t.Run("different models collide", func(t *testing.T) {
		a := shellshape.Normalize("lms load mlx-community/Llama-3-8B")
		b := shellshape.Normalize("lms load TheBloke/Mistral-7B-GGUF")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lms load some-model")
		subshell := shellshape.Normalize("lms load $(dangerous)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
