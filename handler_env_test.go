package shellshape

import "testing"

func TestEnv(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare env", "env", "env"},
		{"env with command", "env bash", "env bash"},
		{"env with assignment and command", "env FOO=bar bash", "env <assign> bash"},
		{"env with multiple assignments", "env FOO=bar BAZ=qux python script.py", "env <assign> python <path>"},

		// Flags
		{"ignore environment", "env -i bash", "env -i bash"},
		{"unset variable", "env -u HOME bash", "env -u <var> bash"},
		{"unset multiple", "env -u PATH -u HOME bash", "env -u <var> -u <var> bash"},
		{"chdir", "env -C /tmp ls", "env -C <path> ls"},
		{"split-string", "env -S 'FOO=bar bash'", "env -S <str>"},
		{"null terminated", "env -0", "env -0"},
		{"long unset", "env --unset=HOME bash", "env --unset=<var> bash"},
		{"long chdir", "env --chdir=/tmp ls", "env --chdir=<path> ls"},

		// Combined flags and assignments
		{"flag then assignment", "env -i FOO=bar bash", "env -i <assign> bash"},
		{"unset then assignment", "env -u HOME FOO=bar bash", "env -u <var> <assign> bash"},

		// Command with arguments
		{"command with args", "env FOO=bar node app.js --port 3000", "env <assign> node <path> --port N"},

		// Redirect
		{"with redirect", "env FOO=bar bash > /tmp/out", "env <assign> bash > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different assignments collide", func(t *testing.T) {
		a := Normalize("env FOO=bar bash")
		b := Normalize("env DATABASE_URL=postgres://localhost/db bash")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different unset vars collide", func(t *testing.T) {
		a := Normalize("env -u HOME bash")
		b := Normalize("env -u PATH bash")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("env FOO=bar bash")
		subshell := Normalize("env $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
