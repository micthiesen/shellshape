package shellshape

import "testing"

func TestPython(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"script only", "python script.py", "python <path>"},
		{"script with path", "python src/main.py", "python <path>"},
		{"script with args", "python app.py --port 3000 foo", "python <path> <arg>+"},
		{"bare python", "python", "python"},
		{"python3 script", "python3 script.py", "python3 <path>"},

		// Code eval
		{"code short", "python -c 'import sys; print(sys.path)'", "python -c <code>"},
		{"code with args", "python -c 'print(1)' arg1 arg2", "python -c <code> <arg>+"},

		// Module mode
		{"module", "python -m pytest", "python -m pytest"},
		{"module with args", "python -m pytest tests/ -v --tb=short", "python -m pytest <path> -v --tb=<val>"},
		{"pip install", "python3 -m pip install requests", "python3 -m pip install requests"},
		{"http server", "python -m http.server 8000", "python -m http.server N"},
		{"pdb", "python -m pdb script.py", "python -m pdb <path>"},

		// Flags with arguments
		{"warning flag", "python -W ignore script.py", "python -W <val> <path>"},
		{"X flag", "python -X dev app.py", "python -X <val> <path>"},
		{"check hash pycs", "python --check-hash-based-pycs always script.py", "python --check-hash-based-pycs <val> <path>"},

		// Boolean flags
		{"unbuffered", "python -u script.py", "python -u <path>"},
		{"no pyc", "python -B script.py", "python -B <path>"},
		{"optimize", "python -O script.py", "python -O <path>"},
		{"interactive", "python -i script.py", "python -i <path>"},
		{"quiet", "python -q script.py", "python -q <path>"},

		// Multiple flags
		{"multi flags", "python -u -B -O script.py arg1", "python -u -B -O <path> <arg>"},
		{"W and script", "python -W error -B app.py", "python -W <val> -B <path>"},

		// Double dash
		{"double dash", "python script.py -- --not-a-flag", "python <path> -- <arg>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different scripts collide", func(t *testing.T) {
		a := Normalize("python server.py --port 3000")
		b := Normalize("python worker.py --port 8080")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different module args collide", func(t *testing.T) {
		a := Normalize("python -m pytest tests/unit")
		b := Normalize("python -m pytest tests/integration")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different code collides", func(t *testing.T) {
		a := Normalize("python -c 'print(1)'")
		b := Normalize("python -c 'import os'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("python script.py")
		subshell := Normalize("python $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
