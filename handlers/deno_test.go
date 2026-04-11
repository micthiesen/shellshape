package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDeno(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"run script", "deno run app.ts", "deno run <path>"},
		{"run with path", "deno run src/server.ts", "deno run <path>"},
		{"run url", "deno run https://deno.land/std/examples/welcome.ts", "deno run <https-uri>"},

		// Permission flags (boolean when no value, value when = syntax)
		{"allow-net boolean", "deno run --allow-net server.ts", "deno run --allow-net <path>"},
		{"allow-net with value", "deno run --allow-net=example.com server.ts", "deno run --allow-net=<val> <path>"},
		{"allow-read with value", "deno run --allow-read=/tmp app.ts", "deno run --allow-read=<val> <path>"},
		{"allow-write boolean", "deno run --allow-write app.ts", "deno run --allow-write <path>"},
		{"deny-net with value", "deno run --deny-net=evil.com app.ts", "deno run --deny-net=<val> <path>"},
		{"allow-all", "deno run -A server.ts", "deno run -A <path>"},
		{"multiple permissions", "deno run --allow-net --allow-read --allow-env app.ts", "deno run --allow-net --allow-read --allow-env <path>"},

		// Path flags
		{"config flag", "deno run --config deno.json app.ts", "deno run --config <path>+"},
		{"config short", "deno run -c deno.json app.ts", "deno run -c <path>+"},
		{"import-map", "deno run --import-map import_map.json app.ts", "deno run --import-map <path>+"},
		{"lock flag", "deno run --lock deno.lock app.ts", "deno run --lock <path>+"},
		{"cert flag", "deno run --cert ca.pem app.ts", "deno run --cert <path>+"},
		{"env-file flag", "deno run --env-file .env app.ts", "deno run --env-file <path>+"},

		// Compile subcommand
		{"compile with output", "deno compile --output myapp app.ts", "deno compile --output <path>+"},
		{"compile with output short", "deno compile -o myapp app.ts", "deno compile -o <path>+"},
		{"compile with target", "deno compile --target x86_64-unknown-linux-gnu app.ts", "deno compile --target <val> <path>"},

		// Test subcommand
		{"test bare", "deno test", "deno test"},
		{"test with path", "deno test src/", "deno test <path>"},
		{"test with filter", "deno test --filter my_test src/", "deno test --filter <val> <path>"},
		{"test with jobs", "deno test --parallel src/", "deno test --parallel <path>"},

		// Fmt and lint
		{"fmt bare", "deno fmt", "deno fmt"},
		{"fmt with files", "deno fmt src/main.ts src/lib.ts", "deno fmt <path>+"},
		{"lint bare", "deno lint", "deno lint"},
		{"lint with config", "deno lint --config deno.json src/", "deno lint --config <path>+"},

		// Task subcommand (task names are structural)
		{"task name", "deno task dev", "deno task dev"},
		{"task with args", "deno task build --mode production", "deno task build --mode production"},

		// Eval subcommand
		{"eval code", "deno eval 'console.log(42)'", "deno eval <code>"},
		{"eval with args", "deno eval 'Deno.exit(0)' arg1 arg2", "deno eval <code> <arg>+"},

		// Install subcommand
		{"install package", "deno install npm:express", "deno install npm:express"},

		// Add/remove
		{"add package", "deno add @std/path", "deno add <path>"},
		{"remove package", "deno remove @std/path", "deno remove <path>"},

		// Check subcommand
		{"check file", "deno check src/main.ts", "deno check <path>"},

		// Info subcommand
		{"info file", "deno info src/main.ts", "deno info <path>"},
		{"info url", "deno info https://deno.land/std/http/server.ts", "deno info <https-uri>"},

		// Serve subcommand
		{"serve file", "deno serve --port 8000 server.ts", "deno serve --port N <path>"},

		// Fused flags with =
		{"config equals", "deno run --config=deno.json app.ts", "deno run --config=<path> <path>"},
		{"target equals", "deno compile --target=x86_64-linux app.ts", "deno compile --target=<val> <path>"},

		// Boolean flags
		{"watch flag", "deno run --watch app.ts", "deno run --watch <path>"},
		{"no-check", "deno run --no-check app.ts", "deno run --no-check <path>"},
		{"unstable", "deno run --unstable app.ts", "deno run --unstable <path>"},
		{"quiet", "deno run -q app.ts", "deno run -q <path>"},

		// Reload flag (boolean or =value)
		{"reload boolean", "deno run --reload app.ts", "deno run --reload <path>"},
		{"reload with value", "deno run --reload=https://deno.land/std app.ts", "deno run --reload=<val> <path>"},

		// Bare deno (repl)
		{"bare deno", "deno", "deno"},

		// Init
		{"init bare", "deno init", "deno init"},
		{"init with dir", "deno init my_project", "deno init my_project"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different scripts collide", func(t *testing.T) {
		a := shellshape.Normalize("deno run server.ts")
		b := shellshape.Normalize("deno run client.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different eval code collides", func(t *testing.T) {
		a := shellshape.Normalize("deno eval 'console.log(1)'")
		b := shellshape.Normalize("deno eval 'Deno.exit(0)'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different config paths collide", func(t *testing.T) {
		a := shellshape.Normalize("deno run --config deno.json app.ts")
		b := shellshape.Normalize("deno run --config tsconfig.json app.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("deno run app.ts")
		subshell := shellshape.Normalize("deno run $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
