package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPnpm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// add/remove: package names collapse to <pkg>
		{"add single", "pnpm add react", "pnpm add <pkg>"},
		{"add multiple", "pnpm add react lodash axios", "pnpm add <pkg>+"},
		{"add save-dev", "pnpm add -D typescript @types/node", "pnpm add -D <pkg>+"},
		{"add save-exact", "pnpm add -E react@18.2.0", "pnpm add -E <pkg>"},
		{"add with filter", "pnpm add --filter @myorg/core react", "pnpm add --filter <val> <pkg>"},
		{"remove single", "pnpm remove lodash", "pnpm remove <pkg>"},
		{"remove multiple", "pnpm remove lodash axios", "pnpm remove <pkg>+"},

		// install: no package positionals expected
		{"install bare", "pnpm install", "pnpm install"},
		{"install frozen", "pnpm install --frozen-lockfile", "pnpm install --frozen-lockfile"},
		{"install prod", "pnpm install --prod", "pnpm install --prod"},
		{"install recursive", "pnpm install -r", "pnpm install -r"},
		{"install with config", "pnpm install --config /tmp/pnpm.conf", "pnpm install --config <path>"},
		{"i alias", "pnpm i --frozen-lockfile", "pnpm i --frozen-lockfile"},

		// run: script name verbatim, rest are args
		{"run script", "pnpm run build", "pnpm run build"},
		{"run script with args", "pnpm run test -- --watch --coverage", "pnpm run test -- --watch --coverage"},
		{"run with filter", "pnpm run --filter @myorg/ui build", "pnpm run --filter <val> build"},

		// exec/dlx: command name verbatim, rest are args
		{"exec command", "pnpm exec vitest", "pnpm exec vitest"},
		{"exec with args", "pnpm exec vitest run --reporter json", "pnpm exec vitest run --reporter json"},
		{"dlx command", "pnpm dlx create-react-app my-app", "pnpm dlx create-react-app <arg>"},
		{"dlx with flags", "pnpm dlx --package typescript tsc --init", "pnpm dlx --package <val> tsc --init"},

		// -C directory flag
		{"dir flag short", "pnpm -C /projects/myapp install", "pnpm -C <path> install"},
		// -F filter alias
		{"filter short", "pnpm add -F @myorg/core react", "pnpm add -F <val> <pkg>"},

		// other subcommands: positionals use classifyToken
		{"why package", "pnpm why react", "pnpm why react"},
		{"init bare", "pnpm init", "pnpm init"},
		{"publish", "pnpm publish --no-git-checks", "pnpm publish --no-git-checks"},

		// npm alias
		{"npm install", "npm install", "npm install"},
		{"npm add", "npm add react lodash", "npm add <pkg>+"},
		{"npm run", "npm run build", "npm run build"},
		{"npm exec", "npm exec vitest", "npm exec vitest"},

		// yarn alias
		{"yarn install", "yarn install", "yarn install"},
		{"yarn add", "yarn add react", "yarn add <pkg>"},
		{"yarn remove", "yarn remove lodash", "yarn remove <pkg>"},
		{"yarn run", "yarn run build", "yarn run build"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different packages collide", func(t *testing.T) {
		a := shellshape.Normalize("pnpm add react express lodash")
		b := shellshape.Normalize("pnpm add axios moment dayjs")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different package versions collide", func(t *testing.T) {
		a := shellshape.Normalize("pnpm add react@18.2.0")
		b := shellshape.Normalize("pnpm add react@17.0.0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pnpm add react")
		subshell := shellshape.Normalize("pnpm add $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
