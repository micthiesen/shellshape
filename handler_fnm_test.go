package shellshape

import "testing"

func TestFnm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install version", "fnm install 18", "fnm install <version>"},
		{"install semver", "fnm install 18.17.0", "fnm install <version>"},
		{"install lts", "fnm install --lts", "fnm install --lts"},
		{"install lts named", "fnm install lts/hydrogen", "fnm install <version>"},
		{"install latest", "fnm install --latest", "fnm install --latest"},
		{"install with mirror", "fnm install --node-dist-mirror https://my-mirror.example.com/dist 20", "fnm install --node-dist-mirror <val> <version>"},
		{"install with fnm-dir", "fnm install --fnm-dir /home/user/.fnm 18", "fnm install --fnm-dir <path> <version>"},
		{"install with corepack", "fnm install --corepack-enabled 20", "fnm install --corepack-enabled <version>"},

		// use subcommand
		{"use version", "fnm use 18", "fnm use <version>"},
		{"use semver", "fnm use 20.11.0", "fnm use <version>"},
		{"use no args", "fnm use", "fnm use"},

		// default subcommand
		{"default version", "fnm default 18", "fnm default <version>"},
		{"default semver", "fnm default 20.11.0", "fnm default <version>"},

		// uninstall subcommand
		{"uninstall version", "fnm uninstall 16.0.0", "fnm uninstall <version>"},
		{"uninstall semver", "fnm uninstall 18.17.1", "fnm uninstall <version>"},

		// alias subcommand
		{"alias version to name", "fnm alias 18 my-project", "fnm alias <version> my-project"},
		{"alias semver to name", "fnm alias 20.11.0 production", "fnm alias <version> production"},

		// unalias subcommand
		{"unalias name", "fnm unalias my-project", "fnm unalias my-project"},

		// list/ls subcommand
		{"list", "fnm list", "fnm list"},
		{"ls alias", "fnm ls", "fnm ls"},

		// current subcommand
		{"current", "fnm current", "fnm current"},

		// env subcommand
		{"env bare", "fnm env", "fnm env"},
		{"env with shell", "fnm env --shell bash", "fnm env --shell bash"},
		{"env with json", "fnm env --json", "fnm env --json"},

		// completions subcommand
		{"completions bash", "fnm completions --shell bash", "fnm completions --shell bash"},

		// exec subcommand
		{"exec with using", "fnm exec --using 18 node --version", "fnm exec --using <version> node --version"},
		{"exec with using semver", "fnm exec --using 20.11.0 node -e 'console.log(1)'", "fnm exec --using <version> node -e console.log(1)"},

		// global flags
		{"log-level flag", "fnm install --log-level quiet 18", "fnm install --log-level <val> <version>"},
		{"arch flag", "fnm install --arch x64 18", "fnm install --arch <val> <version>"},
		{"version-file-strategy", "fnm use --version-file-strategy recursive", "fnm use --version-file-strategy <val>"},

		// list-remote
		{"list-remote", "fnm list-remote", "fnm list-remote"},

		// redirects
		{"env with redirect", "fnm env > /tmp/fnm-env.sh", "fnm env > <path>"},
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
	t.Run("different versions collide", func(t *testing.T) {
		a := Normalize("fnm install 18")
		b := Normalize("fnm install 20.11.0")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different fnm-dir paths collide", func(t *testing.T) {
		a := Normalize("fnm install --fnm-dir /home/alice/.fnm 18")
		b := Normalize("fnm install --fnm-dir /home/bob/.fnm 20")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("fnm install 18")
		subshell := Normalize("fnm install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
