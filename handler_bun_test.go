package shellshape

import "testing"

func TestBun(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// bun run
		{"run script name", "bun run dev", "bun run dev"},
		{"run script name lint", "bun run lint", "bun run lint"},
		{"run file path", "bun run ./src/index.ts", "bun run <path>"},
		{"run with watch", "bun run --watch dev", "bun run --watch dev"},
		{"run with hot", "bun run --hot ./server.ts", "bun run --hot <path>"},
		{"run with env-file", "bun run --env-file .env.local dev", "bun run --env-file <path> dev"},
		{"run with config", "bun run --config ./bunfig.toml dev", "bun run --config <path> dev"},

		// bun test
		{"test bare", "bun test", "bun test"},
		{"test with timeout", "bun test --timeout 5000", "bun test --timeout N"},
		{"test with bail", "bun test --bail 1", "bun test --bail N"},
		{"test with file", "bun test src/utils.test.ts", "bun test <path>"},
		{"test combined", "bun test --bail 1 --timeout 5000 src/utils.test.ts", "bun test --bail N --timeout N <path>"},

		// bun install
		{"install bare", "bun install", "bun install"},
		{"install frozen lockfile", "bun install --frozen-lockfile", "bun install --frozen-lockfile"},
		{"install production", "bun install --production", "bun install --production"},

		// bun build
		{"build entry file", "bun build src/index.ts", "bun build <path>"},
		{"build with outdir", "bun build --outdir ./dist src/index.ts", "bun build --outdir <path>+"},
		{"build with target", "bun build --target browser src/index.ts", "bun build --target <val> <path>"},
		{"build with format", "bun build --format esm src/index.ts", "bun build --format <val> <path>"},
		{"build complex", "bun build --outdir ./dist --target node --format esm src/index.ts src/worker.ts", "bun build --outdir <path> --target <val> --format <val> <path>+"},

		// bun x (bunx)
		{"x package", "bun x prisma", "bun x prisma"},
		{"x package with args", "bun x prisma migrate dev", "bun x prisma migrate dev"},
		{"x with scoped package", "bun x @biomejs/biome check", "bun x <path> check"},

		// Flags with =
		{"long flag with value", "bun run --port=3000 dev", "bun run --port=<val> dev"},
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
	t.Run("different file paths collide", func(t *testing.T) {
		a := Normalize("bun run ./src/server.ts")
		b := Normalize("bun run ./src/client.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different timeout values collide", func(t *testing.T) {
		a := Normalize("bun test --timeout 5000")
		b := Normalize("bun test --timeout 10000")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different env-file paths collide", func(t *testing.T) {
		a := Normalize("bun run --env-file .env.local dev")
		b := Normalize("bun run --env-file .env.production dev")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("bun run script-name")
		subshell := Normalize("bun run $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
