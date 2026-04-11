package shellshape

import "testing"

func TestSnap(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install package", "snap install firefox", "snap install firefox"},
		{"install classic", "snap install --classic code", "snap install --classic code"},
		{"install devmode", "snap install --devmode mysnap", "snap install --devmode mysnap"},
		{"install dangerous", "snap install --dangerous foo.snap", "snap install --dangerous foo.snap"},
		{"install channel", "snap install --channel=beta vlc", "snap install --channel=<val> vlc"},
		{"install channel separate", "snap install --channel beta vlc", "snap install --channel <val> vlc"},

		// remove subcommand
		{"remove package", "snap remove firefox", "snap remove firefox"},
		{"remove purge", "snap remove --purge firefox", "snap remove --purge firefox"},

		// refresh subcommand
		{"refresh all", "snap refresh", "snap refresh"},
		{"refresh package", "snap refresh firefox", "snap refresh firefox"},
		{"refresh channel", "snap refresh --channel=stable firefox", "snap refresh --channel=<val> firefox"},

		// revert subcommand
		{"revert package", "snap revert core", "snap revert core"},
		{"revert revision", "snap revert core --revision 5", "snap revert core --revision N"},

		// find subcommand
		{"find query", "snap find image-editor", "snap find <query>"},
		{"find private", "snap find --private", "snap find --private"},

		// info subcommand
		{"info package", "snap info firefox", "snap info firefox"},

		// list subcommand
		{"list all", "snap list", "snap list"},
		{"list all revisions", "snap list --all", "snap list --all"},
		{"list specific", "snap list firefox", "snap list firefox"},

		// enable/disable
		{"enable package", "snap enable foo", "snap enable foo"},
		{"disable package", "snap disable foo", "snap disable foo"},

		// connect/disconnect
		{"connect plug", "snap connect foo:camera :camera", "snap connect foo:camera :camera"},
		{"disconnect plug", "snap disconnect foo:camera", "snap disconnect foo:camera"},

		// set/get
		{"set key value", "snap set foo bar=10", "snap set foo <val>"},
		{"set multiple", "snap set foo bar=10 baz=hello", "snap set foo <val>+"},
		{"get key", "snap get foo bar", "snap get foo bar"},

		// change/watch/abort with numeric IDs
		{"change id", "snap change 123", "snap change N"},
		{"watch id", "snap watch 456", "snap watch N"},
		{"abort id", "snap abort 789", "snap abort N"},

		// changes (no positional)
		{"changes", "snap changes", "snap changes"},

		// download
		{"download package", "snap download core", "snap download core"},

		// ack (file path)
		{"ack assertion", "snap ack foo.assert", "snap ack <dotted-id>"},

		// login/logout
		{"login", "snap login", "snap login"},
		{"logout", "snap logout", "snap logout"},

		// interfaces
		{"interfaces", "snap interfaces", "snap interfaces"},

		// version
		{"version", "snap version", "snap version"},

		// Subshell as val flag argument (--channel)
		{"channel subshell", "snap install --channel $(get-channel) vlc", "snap install --channel $(get-channel) vlc"},
		// Subshell as num flag argument (--revision)
		{"revision subshell", "snap revert core --revision $(get-rev)", "snap revert core --revision $(get-rev)"},
		// Unknown subcommand falls back to generic classification
		{"unknown subcommand", "snap run myapp", "snap run myapp"},
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
	t.Run("different find queries collide", func(t *testing.T) {
		a := Normalize("snap find postgres")
		b := Normalize("snap find redis")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different change IDs collide", func(t *testing.T) {
		a := Normalize("snap change 100")
		b := Normalize("snap change 999")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("snap install literal-arg")
		subshell := Normalize("snap install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
