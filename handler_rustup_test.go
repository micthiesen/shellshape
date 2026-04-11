package shellshape

import "testing"

func TestRustup(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install stable", "rustup install stable", "rustup install <toolchain>"},
		{"install nightly", "rustup install nightly", "rustup install <toolchain>"},
		{"install specific version", "rustup install 1.70.0", "rustup install <toolchain>"},
		{"install nightly date", "rustup install nightly-2024-01-15", "rustup install <toolchain>"},
		{"install with profile", "rustup install --profile minimal stable", "rustup install --profile <val> <toolchain>"},

		// update subcommand
		{"update bare", "rustup update", "rustup update"},
		{"update specific", "rustup update nightly", "rustup update <toolchain>"},
		{"update no-self-update", "rustup update --no-self-update", "rustup update --no-self-update"},

		// default subcommand
		{"default stable", "rustup default stable", "rustup default <toolchain>"},
		{"default nightly", "rustup default nightly", "rustup default <toolchain>"},

		// show subcommand
		{"show bare", "rustup show", "rustup show"},
		{"show active-toolchain", "rustup show active-toolchain", "rustup show active-toolchain"},

		// run subcommand - first positional is toolchain, rest are verbatim command
		{"run cargo build", "rustup run nightly cargo build", "rustup run <toolchain> cargo build"},
		{"run cargo test", "rustup run stable cargo test --release", "rustup run <toolchain> cargo test --release"},
		{"run rustc version", "rustup run 1.70.0 rustc --version", "rustup run <toolchain> rustc --version"},

		// toolchain subcommand (has second subcommand)
		{"toolchain list", "rustup toolchain list", "rustup toolchain list"},
		{"toolchain install", "rustup toolchain install nightly", "rustup toolchain install <toolchain>"},
		{"toolchain remove", "rustup toolchain remove nightly-2024-01-01", "rustup toolchain remove <toolchain>"},
		{"toolchain link", "rustup toolchain link my-toolchain /path/to/toolchain", "rustup toolchain link <toolchain> <path>"},

		// target subcommand (has second subcommand)
		{"target list", "rustup target list", "rustup target list"},
		{"target add", "rustup target add wasm32-unknown-unknown", "rustup target add <target>"},
		{"target add with toolchain", "rustup target add x86_64-pc-windows-msvc --toolchain nightly", "rustup target add <target> --toolchain <toolchain>"},
		{"target remove", "rustup target remove wasm32-wasi", "rustup target remove <target>"},

		// component subcommand (has second subcommand)
		{"component list", "rustup component list", "rustup component list"},
		{"component add", "rustup component add rustfmt", "rustup component add <component>"},
		{"component add with toolchain", "rustup component add clippy --toolchain nightly", "rustup component add <component> --toolchain <toolchain>"},
		{"component remove", "rustup component remove rust-docs", "rustup component remove <component>"},

		// override subcommand (has second subcommand)
		{"override list", "rustup override list", "rustup override list"},
		{"override set", "rustup override set nightly", "rustup override set <toolchain>"},
		{"override set with path", "rustup override set nightly --path /home/user/project", "rustup override set <toolchain> --path <path>"},
		{"override unset", "rustup override unset", "rustup override unset"},

		// self subcommand (has second subcommand, all verbatim)
		{"self update", "rustup self update", "rustup self update"},
		{"self uninstall", "rustup self uninstall", "rustup self uninstall"},

		// which subcommand
		{"which rustc", "rustup which rustc", "rustup which <component>"},
		{"which cargo", "rustup which cargo", "rustup which <component>"},

		// doc subcommand
		{"doc bare", "rustup doc", "rustup doc"},
		{"doc with flag", "rustup doc --std", "rustup doc --std"},

		// doctor subcommand
		{"doctor", "rustup doctor", "rustup doctor"},

		// completions subcommand
		{"completions bash", "rustup completions bash cargo", "rustup completions bash cargo"},

		// global flags (after subcommand)
		{"quiet flag", "rustup update --quiet", "rustup update --quiet"},

		// redirects
		{"with redirect", "rustup show > /tmp/toolchains.txt", "rustup show > <path>"},
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
	t.Run("different toolchains collide", func(t *testing.T) {
		a := Normalize("rustup install stable")
		b := Normalize("rustup install nightly-2024-01-15")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different targets collide", func(t *testing.T) {
		a := Normalize("rustup target add wasm32-unknown-unknown")
		b := Normalize("rustup target add x86_64-pc-windows-msvc")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("rustup install stable")
		subshell := Normalize("rustup install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
