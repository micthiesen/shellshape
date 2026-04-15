package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestBrew(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install single formula", "brew install wget", "brew install <pkg>"},
		{"install multiple formulas", "brew install node python go", "brew install <pkg>+"},
		{"install cask", "brew install --cask firefox", "brew install --cask <pkg>"},
		{"install HEAD", "brew install --HEAD neovim", "brew install --HEAD <pkg>"},
		{"install force verbose", "brew install --force --verbose node", "brew install --force --verbose <pkg>"},
		{"install formula flag", "brew install --formula gcc", "brew install --formula <pkg>"},

		// uninstall subcommand
		{"uninstall formula", "brew uninstall node", "brew uninstall <pkg>"},
		{"uninstall cask", "brew uninstall --cask firefox", "brew uninstall --cask <pkg>"},
		{"remove alias", "brew remove wget", "brew remove <pkg>"},

		// upgrade subcommand
		{"upgrade all", "brew upgrade", "brew upgrade"},
		{"upgrade specific", "brew upgrade node", "brew upgrade <pkg>"},
		{"upgrade multiple", "brew upgrade node python go", "brew upgrade <pkg>+"},
		{"upgrade greedy cask", "brew upgrade --greedy --cask", "brew upgrade --greedy --cask"},

		// info subcommand
		{"info formula", "brew info node", "brew info <pkg>"},
		{"info json", "brew info --json=v2 node", "brew info --json=<val> <pkg>"},

		// search subcommand
		{"search text", "brew search postgres", "brew search <query>"},
		{"search with flag", "brew search --cask font", "brew search --cask <query>"},

		// list subcommand (formula names stay verbatim; list is a query-ish view)
		{"list all", "brew list", "brew list"},
		{"list cask", "brew list --cask", "brew list --cask"},
		{"list specific", "brew list node", "brew list node"},

		// tap/untap subcommands
		{"tap repo", "brew tap homebrew/cask-fonts", "brew tap homebrew/cask-fonts"},
		{"tap list", "brew tap", "brew tap"},
		{"untap repo", "brew untap homebrew/cask-fonts", "brew untap homebrew/cask-fonts"},

		// services subcommand: action verbatim, service name <pkg>
		{"services list", "brew services list", "brew services list"},
		{"services start", "brew services start postgresql", "brew services start <pkg>"},
		{"services stop", "brew services stop redis", "brew services stop <pkg>"},
		{"services restart", "brew services restart nginx", "brew services restart <pkg>"},

		// deps subcommand
		{"deps formula", "brew deps node", "brew deps <pkg>"},
		{"deps tree", "brew deps --tree node", "brew deps --tree <pkg>"},

		// reinstall
		{"reinstall formula", "brew reinstall node", "brew reinstall <pkg>"},

		// pin/unpin
		{"pin formula", "brew pin node", "brew pin <pkg>"},
		{"unpin formula", "brew unpin node", "brew unpin <pkg>"},

		// boolean flags preserved
		{"quiet flag", "brew install -q node", "brew install -q <pkg>"},
		{"verbose short", "brew install -v wget", "brew install -v <pkg>"},
		{"dry-run", "brew upgrade --dry-run", "brew upgrade --dry-run"},

		// other/unknown subcommands fall back to classifyToken
		{"config", "brew config", "brew config"},
		{"doctor", "brew doctor", "brew doctor"},
		{"cleanup with path", "brew cleanup --prune=7", "brew cleanup --prune=<val>"},
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
	t.Run("different search terms collide", func(t *testing.T) {
		a := shellshape.Normalize("brew search postgres")
		b := shellshape.Normalize("brew search redis")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different formula names collide", func(t *testing.T) {
		a := shellshape.Normalize("brew install sleepwatcher")
		b := shellshape.Normalize("brew install yabai")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different service names collide", func(t *testing.T) {
		a := shellshape.Normalize("brew services start sleepwatcher")
		b := shellshape.Normalize("brew services start nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("brew install literal-arg")
		subshell := shellshape.Normalize("brew install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
