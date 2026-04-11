package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDotenvx(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare dotenvx run", "dotenvx run -- npm start", "dotenvx run -- npm start"},
		{"bare dotenv", "dotenv -- npm start", "dotenv -- npm start"},

		// Flags with path arguments
		{"-f env file", "dotenvx run -f .env.production -- node index.js", "dotenvx run -f <path> -- node <path>"},
		{"--env-file", "dotenvx run --env-file .env.local -- node index.js", "dotenvx run --env-file <path> -- node <path>"},
		{"multiple -f", "dotenvx run -f .env.local -f .env -- node index.js", "dotenvx run -f <path> -f <path> -- node <path>"},
		{"-fv vault file", "dotenvx run -fv .env.vault -- node index.js", "dotenvx run -fv <path> -- node <path>"},
		{"--env-vault-file", "dotenvx run --env-vault-file .env.vault -- node index.js", "dotenvx run --env-vault-file <path> -- node <path>"},

		// Flags with value arguments
		{"-e env var", "dotenvx run -e HELLO=World -- node index.js", "dotenvx run -e <val> -- node <path>"},
		{"--env inline", "dotenvx run --env HELLO=World -- node index.js", "dotenvx run --env <val> -- node <path>"},
		{"--convention", "dotenvx run --convention nextjs -- npm run dev", "dotenvx run --convention <val> -- npm run dev"},

		// Boolean flags
		{"--overload", "dotenvx run --overload -f .env -- node index.js", "dotenvx run --overload -f <path> -- node <path>"},
		{"-o overload", "dotenvx run -o -f .env -- node index.js", "dotenvx run -o -f <path> -- node <path>"},

		// dotenv-cli alias (-e takes a value; for dotenv-cli it's a path, for dotenvx it's KEY=VAL)
		{"dotenv -e", "dotenv -e .env2 -- mvn exec:java", "dotenv -e <val> -- mvn exec:java"},
		{"dotenv multiple -e", "dotenv -e .env3 -e .env4 -- node app.js", "dotenv -e <val> -e <val> -- node <path>"},
		{"dotenv --override", "dotenv --override -e .env -- node app.js", "dotenv --override -e <val> -- node <path>"},
		{"dotenv -o", "dotenv -o -- npm test", "dotenv -o -- npm test"},
		{"dotenv --no-expand", "dotenv --no-expand -- node app.js", "dotenv --no-expand -- node <path>"},

		// Command after -- gets generic classification
		{"command with paths after --", "dotenvx run -- python /usr/local/bin/app.py", "dotenvx run -- python <path>"},
		{"command with flags after --", "dotenvx run -- flask --app index run", "dotenvx run -- flask --app index run"},

		// Redirect
		{"with redirect", "dotenvx run -- node index.js > /tmp/out.log", "dotenvx run -- node <path> > <path>"},

		// No -- separator (everything is flags/positionals for dotenvx)
		{"no separator", "dotenvx run -f .env node index.js", "dotenvx run -f <path> node <path>"},
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
	t.Run("different env files collide", func(t *testing.T) {
		a := shellshape.Normalize("dotenvx run -f .env.production -- node index.js")
		b := shellshape.Normalize("dotenvx run -f .env.staging -- node index.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different env values collide", func(t *testing.T) {
		a := shellshape.Normalize("dotenvx run -e HELLO=World -- node index.js")
		b := shellshape.Normalize("dotenvx run -e DB_HOST=localhost -- node index.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different commands after -- collide on paths", func(t *testing.T) {
		a := shellshape.Normalize("dotenvx run -- node /app/server.js")
		b := shellshape.Normalize("dotenvx run -- node /app/client.js")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("dotenvx run -f .env -- node index.js")
		subshell := shellshape.Normalize("dotenvx run -f $(dangerous-command) -- node index.js")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
