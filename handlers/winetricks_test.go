package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWinetricks(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"single verb", "winetricks dxvk", "winetricks dxvk"},
		{"multiple verbs", "winetricks dxvk vcrun2019 dotnet48", "winetricks dxvk vcrun2019 dotnet48"},
		{"unattended", "winetricks --unattended dxvk", "winetricks --unattended dxvk"},
		{"short quiet", "winetricks -q dxvk vcrun2019", "winetricks -q dxvk vcrun2019"},
		{"force flag", "winetricks --force dxvk", "winetricks --force dxvk"},
		{"gui flag", "winetricks --gui", "winetricks --gui"},
		{"prefix flag", "winetricks --prefix /home/user/.wine dxvk", "winetricks --prefix <path> dxvk"},
		{"prefix with quiet", "winetricks -q --prefix /opt/wine-prefix vcrun2019", "winetricks -q --prefix <path> vcrun2019"},
		{"multiple flags and verbs", "winetricks --unattended --force dxvk vcrun2019", "winetricks --unattended --force dxvk vcrun2019"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different prefixes collide", func(t *testing.T) {
		a := shellshape.Normalize("winetricks --prefix /home/user/.wine1 dxvk")
		b := shellshape.Normalize("winetricks --prefix /home/user/.wine2 dxvk")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("winetricks dxvk")
		subshell := shellshape.Normalize("winetricks $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
