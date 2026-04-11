package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestNix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands (no extra args)
		{"repl", "nix repl", "nix repl"},
		{"store gc", "nix store gc", "nix store gc"},
		{"flake update", "nix flake update", "nix flake update"},
		{"doctor", "nix doctor", "nix doctor"},

		// build subcommand
		{"build flake ref", "nix build nixpkgs#hello", "nix build <pkg>"},
		{"build local flake", "nix build .#mypackage", "nix build <pkg>"},
		{"build current dir", "nix build", "nix build"},
		{"build with no-link", "nix build --no-link nixpkgs#hello", "nix build --no-link <pkg>"},
		{"build with out-link", "nix build -o result-bin nixpkgs#hello", "nix build -o <path> <pkg>"},
		{"build dry-run", "nix build --dry-run .#default", "nix build --dry-run <pkg>"},
		{"build with file", "nix build -f default.nix", "nix build -f <path>"},
		{"build json", "nix build --json nixpkgs#hello", "nix build --json <pkg>"},

		// shell subcommand
		{"shell single pkg", "nix shell nixpkgs#gcc", "nix shell <pkg>"},
		{"shell multiple pkgs", "nix shell nixpkgs#gcc nixpkgs#cmake nixpkgs#ninja", "nix shell <pkg>+"},

		// run subcommand
		{"run simple", "nix run nixpkgs#hello", "nix run <pkg>"},
		{"run with args", "nix run nixpkgs#hello -- --greeting hi", "nix run <pkg> -- --greeting hi"},

		// develop subcommand
		{"develop default", "nix develop", "nix develop"},
		{"develop impure", "nix develop --impure", "nix develop --impure"},
		{"develop with flake", "nix develop github:user/repo", "nix develop <pkg>"},

		// search subcommand
		{"search simple", "nix search nixpkgs hello", "nix search <pkg> <str>"},
		{"search multi terms", "nix search nixpkgs python web", "nix search <pkg> <str>+"},

		// eval subcommand
		{"eval with expr", "nix eval --expr '1 + 1'", "nix eval --expr <expr>"},
		{"eval flake attr", "nix eval .#packages.x86_64-linux.default", "nix eval <pkg>"},

		// profile subcommand
		{"profile install", "nix profile install nixpkgs#hello", "nix profile install <pkg>"},
		{"profile install multi", "nix profile install nixpkgs#hello nixpkgs#cowsay", "nix profile install <pkg>+"},
		{"profile remove", "nix profile remove 0", "nix profile remove N"},
		{"profile list", "nix profile list", "nix profile list"},

		// flake subcommand
		{"flake init", "nix flake init", "nix flake init"},
		{"flake show", "nix flake show", "nix flake show"},
		{"flake show remote", "nix flake show github:NixOS/nixpkgs", "nix flake show <pkg>"},
		{"flake lock update-input", "nix flake lock --update-input nixpkgs", "nix flake lock --update-input <val>"},
		{"flake check", "nix flake check", "nix flake check"},

		// Flags with arguments
		{"override-input", "nix build --override-input nixpkgs github:NixOS/nixpkgs/master .#foo", "nix build --override-input <val>+ <pkg>"},
		{"inputs-from", "nix build --inputs-from . .#foo", "nix build --inputs-from <val> <pkg>"},
		{"include path", "nix build -I nixpkgs=/my/path .#foo", "nix build -I <path> <pkg>"},
		{"log-format", "nix build --log-format bar .#foo", "nix build --log-format <val> <pkg>"},
		{"option flag", "nix build --option sandbox false .#foo", "nix build --option <val>+ <pkg>"},
		{"arg flag", "nix eval --arg name 'expr' --expr 'name'", "nix eval --arg <val>+ --expr <expr>"},
		{"eval-store", "nix build --eval-store /tmp/store .#foo", "nix build --eval-store <val> <pkg>"},

		// Fused --flag=value
		{"fused out-link", "nix build --out-link=result-bin .#foo", "nix build --out-link=<path> <pkg>"},
		{"fused log-format", "nix build --log-format=bar .#foo", "nix build --log-format=<val> <pkg>"},

		// Global boolean flags
		{"verbose", "nix build -v .#hello", "nix build -v <pkg>"},
		{"print-build-logs", "nix build -L .#hello", "nix build -L <pkg>"},
		{"offline", "nix build --offline .#hello", "nix build --offline <pkg>"},
		{"refresh", "nix build --refresh nixpkgs#hello", "nix build --refresh <pkg>"},

		// copy subcommand
		{"copy to store", "nix copy --to ssh://server /nix/store/abc123-hello", "nix copy --to <val> <path>"},

		// hash subcommand
		{"hash path", "nix hash path ./my-source", "nix hash path <path>"},

		// log subcommand
		{"log flake", "nix log nixpkgs#hello", "nix log <pkg>"},

		// why-depends
		{"why-depends", "nix why-depends nixpkgs#hello nixpkgs#glibc", "nix why-depends <pkg>+"},

		// path-info
		{"path-info", "nix path-info --closure-size /nix/store/abc123-hello", "nix path-info --closure-size <path>"},

		// Redirects
		{"build with redirect", "nix build .#foo > output.txt", "nix build <pkg> > <path>"},
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
	t.Run("different flake refs collide", func(t *testing.T) {
		a := shellshape.Normalize("nix build nixpkgs#hello")
		b := shellshape.Normalize("nix build nixpkgs#cowsay")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different local flake refs collide", func(t *testing.T) {
		a := shellshape.Normalize("nix build .#myapp")
		b := shellshape.Normalize("nix build .#otherapp")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different search terms collide", func(t *testing.T) {
		a := shellshape.Normalize("nix search nixpkgs python")
		b := shellshape.Normalize("nix search nixpkgs nodejs")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("nix build literal-arg")
		subshell := shellshape.Normalize("nix build $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
