package shellshape

import "testing"

func TestTsx(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple script", "tsx app.ts", "tsx <script>"},
		{"script with extension tsx", "tsx component.tsx", "tsx <script>"},
		{"script with path", "tsx src/server.ts", "tsx <script>"},

		// Watch subcommand
		{"watch mode", "tsx watch app.ts", "tsx watch <script>"},
		{"watch with path", "tsx watch src/index.ts", "tsx watch <script>"},

		// Flags with path arguments
		{"tsconfig flag", "tsx --tsconfig tsconfig.custom.json app.ts", "tsx --tsconfig <path> <script>"},
		{"tsconfig with watch", "tsx watch --tsconfig tsconfig.json src/server.ts", "tsx watch --tsconfig <path> <script>"},

		// Flags with module arguments
		{"require short", "tsx -r dotenv/config app.ts", "tsx -r <module> <script>"},
		{"require long", "tsx --require dotenv/config app.ts", "tsx --require <module> <script>"},
		{"import flag", "tsx --import ./setup.ts app.ts", "tsx --import <module> <script>"},

		// Boolean flags
		{"inspect flag", "tsx --inspect src/debug.ts", "tsx --inspect <script>"},
		{"inspect-brk flag", "tsx --inspect-brk src/debug.ts", "tsx --inspect-brk <script>"},
		{"no-cache flag", "tsx --no-cache app.ts", "tsx --no-cache <script>"},

		// Script arguments collapsed
		{"script with args", "tsx app.ts arg1 arg2 arg3", "tsx <script> <arg>+"},
		{"watch with script args", "tsx watch server.ts --port 3000 --host localhost", "tsx watch <script> <arg>+"},
		{"script with mixed args", "tsx app.ts foo bar --verbose", "tsx <script> <arg>+"},

		// Combination
		{"full combo", "tsx --tsconfig tsconfig.json -r dotenv/config src/main.ts --env production", "tsx --tsconfig <path> -r <module> <script> <arg>+"},
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
		a := Normalize("tsx src/server.ts")
		b := Normalize("tsx src/worker.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different script args collide", func(t *testing.T) {
		a := Normalize("tsx app.ts --port 3000 --host 0.0.0.0")
		b := Normalize("tsx app.ts --port 8080 --host localhost")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different modules collide", func(t *testing.T) {
		a := Normalize("tsx -r dotenv/config app.ts")
		b := Normalize("tsx -r tsconfig-paths/register app.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tsx app.ts")
		subshell := Normalize("tsx $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
