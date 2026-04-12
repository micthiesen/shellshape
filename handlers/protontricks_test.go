package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestProtontricks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"app id with verb", "protontricks 123456 dxvk", "protontricks N dxvk"},
		{"app id with multiple verbs", "protontricks 123456 dxvk vcrun2019", "protontricks N dxvk vcrun2019"},
		{"no runtime", "protontricks --no-runtime 123456 dxvk", "protontricks --no-runtime N dxvk"},
		{"search", "protontricks -s 'Elden Ring'", "protontricks -s <val>"},
		{"search long", "protontricks --search 'Elden Ring'", "protontricks --search <val>"},
		{"command flag", "protontricks -c 'winecfg' 123456", "protontricks -c <val> N"},
		{"gui flag", "protontricks --gui", "protontricks --gui"},
		{"verbose", "protontricks -v 123456 dxvk", "protontricks -v N dxvk"},
		{"verbose long", "protontricks --verbose 123456 dxvk", "protontricks --verbose N dxvk"},
		{"full example", "protontricks --no-runtime -v 123456 dxvk vcrun2019 dotnet48", "protontricks --no-runtime -v N dxvk vcrun2019 dotnet48"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different app ids collide", func(t *testing.T) {
		a := shellshape.Normalize("protontricks 123456 dxvk")
		b := shellshape.Normalize("protontricks 789012 dxvk")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("protontricks 123456")
		subshell := shellshape.Normalize("protontricks $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
