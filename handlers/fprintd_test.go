package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFprintd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// fprintd-enroll
		{"enroll bare", "fprintd-enroll", "fprintd-enroll"},
		{"enroll user", "fprintd-enroll michael", "fprintd-enroll <val>"},
		{"enroll finger", "fprintd-enroll -f right-index-finger", "fprintd-enroll -f right-index-finger"},
		{"enroll finger user", "fprintd-enroll -f left-thumb michael", "fprintd-enroll -f left-thumb <val>"},

		// fprintd-verify
		{"verify bare", "fprintd-verify", "fprintd-verify"},
		{"verify user", "fprintd-verify michael", "fprintd-verify <val>"},
		{"verify finger", "fprintd-verify -f right-index-finger michael", "fprintd-verify -f right-index-finger <val>"},

		// fprintd-list
		{"list user", "fprintd-list michael", "fprintd-list <val>"},

		// fprintd-delete
		{"delete user", "fprintd-delete michael", "fprintd-delete <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different usernames collide", func(t *testing.T) {
		a := shellshape.Normalize("fprintd-enroll -f right-index-finger michael")
		b := shellshape.Normalize("fprintd-enroll -f right-index-finger root")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fprintd-enroll michael")
		subshell := shellshape.Normalize("fprintd-enroll $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
