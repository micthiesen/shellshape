package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestCargo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"build", "cargo build", "cargo build"},
		{"build release", "cargo build --release", "cargo build --release"},
		{"run", "cargo run", "cargo run"},
		{"test", "cargo test", "cargo test"},
		{"check", "cargo check", "cargo check"},
		{"clippy", "cargo clippy", "cargo clippy"},
		{"fmt", "cargo fmt", "cargo fmt"},
		{"clean", "cargo clean", "cargo clean"},

		// Package names (structural - kept verbatim)
		{"add package", "cargo add serde", "cargo add serde"},
		{"add with features", "cargo add serde --features derive", "cargo add serde --features <val>"},
		{"remove package", "cargo remove tokio", "cargo remove tokio"},
		{"install package", "cargo install ripgrep", "cargo install ripgrep"},
		{"update package", "cargo update -p rand", "cargo update -p rand"},

		// Structural flags
		{"target triple", "cargo build --target x86_64-unknown-linux-gnu", "cargo build --target x86_64-unknown-linux-gnu"},
		{"bin name", "cargo run --bin mybin", "cargo run --bin mybin"},
		{"example name", "cargo run --example hello", "cargo run --example hello"},
		{"package flag", "cargo test --package mylib", "cargo test --package mylib"},
		{"short package", "cargo build -p mylib", "cargo build -p mylib"},

		// Numeric flags
		{"jobs", "cargo build -j 8", "cargo build -j N"},
		{"jobs long", "cargo build --jobs 4", "cargo build --jobs N"},

		// Path flags
		{"manifest path", "cargo build --manifest-path /path/to/Cargo.toml", "cargo build --manifest-path <path>"},
		{"target dir", "cargo build --target-dir /tmp/build", "cargo build --target-dir <path>"},

		// Value flags
		{"features", "cargo build --features serde,tokio", "cargo build --features <val>"},
		{"color", "cargo build --color always", "cargo build --color <val>"},

		// Test with pass-through args
		{"test filter", "cargo test my_test_name", "cargo test my_test_name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different features collide", func(t *testing.T) {
		a := shellshape.Normalize("cargo build --features serde,tokio")
		b := shellshape.Normalize("cargo build --features reqwest,hyper")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("cargo build --manifest-path /project/a/Cargo.toml")
		b := shellshape.Normalize("cargo build --manifest-path /project/b/Cargo.toml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("cargo build --features serde")
		subshell := shellshape.Normalize("cargo build --features $(dangerous)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
