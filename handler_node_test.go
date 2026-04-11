package shellshape

import "testing"

func TestNode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"script only", "node app.js", "node <path>"},
		{"script with path", "node src/index.ts", "node <path>"},
		{"script with args", "node app.js --port 3000 foo", "node <path> <arg>+"},

		// Code eval
		{"eval short", "node -e 'console.log(1)'", "node -e <code>"},
		{"eval long", "node --eval 'process.exit(0)'", "node --eval <code>"},
		{"eval with args", "node -e 'console.log(1)' arg1 arg2", "node -e <code> <arg>+"},
		{"check flag", "node -c 'var x = 1'", "node -c <code>"},

		// Flags with path arguments
		{"require short", "node -r ts-node/register src/index.ts", "node -r <path>+"},
		{"require long", "node --require dotenv/config app.js", "node --require <path>+"},
		{"loader", "node --loader tsx src/index.ts", "node --loader <path>+"},
		{"experimental-loader", "node --experimental-loader ./loader.mjs app.js", "node --experimental-loader <path>+"},
		{"env-file", "node --env-file .env app.js", "node --env-file <path>+"},
		{"import", "node --import ./setup.js app.js", "node --import <path>+"},

		// Flags with = syntax
		{"inspect equals", "node --inspect=9229 app.js", "node --inspect=<val> <path>"},
		{"inspect-brk equals", "node --inspect-brk=0.0.0.0:9229 app.js", "node --inspect-brk=<val> <path>"},
		{"max-old-space-size", "node --max-old-space-size=4096 app.js", "node --max-old-space-size=<val> <path>"},

		// Boolean flags
		{"inspect no value", "node --inspect app.js", "node --inspect <path>"},
		{"inspect-brk no value", "node --inspect-brk app.js", "node --inspect-brk <path>"},
		{"enable-source-maps", "node --enable-source-maps app.js", "node --enable-source-maps <path>"},

		// Conditions flag
		{"conditions short", "node -C development app.js", "node -C <val> <path>"},
		{"conditions long", "node --conditions development app.js", "node --conditions <val> <path>"},

		// Multiple flags
		{"complex invocation", "node --inspect -r tsconfig-paths/register src/server.ts", "node --inspect -r <path>+"},
		{"env-file and require", "node --env-file .env -r dotenv/config app.js", "node --env-file <path> -r <path>+"},

		// Double dash
		{"double dash", "node app.js -- --not-a-flag foo", "node <path> -- <arg>+"},

		// Bare node (no args)
		{"no args", "node", "node"},
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
		a := Normalize("node server.js --port 3000")
		b := Normalize("node worker.js --port 8080")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different eval code collides", func(t *testing.T) {
		a := Normalize("node -e 'console.log(1)'")
		b := Normalize("node -e 'process.exit(0)'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("node app.js")
		subshell := Normalize("node $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
