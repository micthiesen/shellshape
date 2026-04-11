package shellshape

import "testing"

func TestBrew(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install single formula", "brew install wget", "brew install wget"},
		{"install multiple formulas", "brew install node python go", "brew install node python go"},
		{"install cask", "brew install --cask firefox", "brew install --cask firefox"},
		{"install HEAD", "brew install --HEAD neovim", "brew install --HEAD neovim"},
		{"install force verbose", "brew install --force --verbose node", "brew install --force --verbose node"},
		{"install formula flag", "brew install --formula gcc", "brew install --formula gcc"},

		// uninstall subcommand
		{"uninstall formula", "brew uninstall node", "brew uninstall node"},
		{"uninstall cask", "brew uninstall --cask firefox", "brew uninstall --cask firefox"},
		{"remove alias", "brew remove wget", "brew remove wget"},

		// upgrade subcommand
		{"upgrade all", "brew upgrade", "brew upgrade"},
		{"upgrade specific", "brew upgrade node", "brew upgrade node"},
		{"upgrade multiple", "brew upgrade node python go", "brew upgrade node python go"},
		{"upgrade greedy cask", "brew upgrade --greedy --cask", "brew upgrade --greedy --cask"},

		// info subcommand
		{"info formula", "brew info node", "brew info node"},
		{"info json", "brew info --json=v2 node", "brew info --json=<val> node"},

		// search subcommand
		{"search text", "brew search postgres", "brew search <query>"},
		{"search with flag", "brew search --cask font", "brew search --cask <query>"},

		// list subcommand
		{"list all", "brew list", "brew list"},
		{"list cask", "brew list --cask", "brew list --cask"},
		{"list specific", "brew list node", "brew list node"},

		// tap/untap subcommands
		{"tap repo", "brew tap homebrew/cask-fonts", "brew tap homebrew/cask-fonts"},
		{"tap list", "brew tap", "brew tap"},
		{"untap repo", "brew untap homebrew/cask-fonts", "brew untap homebrew/cask-fonts"},

		// services subcommand
		{"services list", "brew services list", "brew services list"},
		{"services start", "brew services start postgresql", "brew services start postgresql"},
		{"services stop", "brew services stop redis", "brew services stop redis"},
		{"services restart", "brew services restart nginx", "brew services restart nginx"},

		// deps subcommand
		{"deps formula", "brew deps node", "brew deps node"},
		{"deps tree", "brew deps --tree node", "brew deps --tree node"},

		// reinstall
		{"reinstall formula", "brew reinstall node", "brew reinstall node"},

		// pin/unpin
		{"pin formula", "brew pin node", "brew pin node"},
		{"unpin formula", "brew unpin node", "brew unpin node"},

		// boolean flags preserved
		{"quiet flag", "brew install -q node", "brew install -q node"},
		{"verbose short", "brew install -v wget", "brew install -v wget"},
		{"dry-run", "brew upgrade --dry-run", "brew upgrade --dry-run"},

		// other/unknown subcommands fall back to classifyToken
		{"config", "brew config", "brew config"},
		{"doctor", "brew doctor", "brew doctor"},
		{"cleanup with path", "brew cleanup --prune=7", "brew cleanup --prune=<val>"},
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
	t.Run("different search terms collide", func(t *testing.T) {
		a := Normalize("brew search postgres")
		b := Normalize("brew search redis")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("brew install literal-arg")
		subshell := Normalize("brew install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
