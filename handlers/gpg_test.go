package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGpg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"gen key", "gpg --gen-key", "gpg --gen-key"},
		{"list keys", "gpg --list-keys", "gpg --list-keys"},
		{"list secret keys", "gpg --list-secret-keys", "gpg --list-secret-keys"},

		// Encrypt / decrypt
		{"encrypt with recipient", "gpg --encrypt --recipient user@example.com file.txt", "gpg --encrypt --recipient <val> <path>"},
		{"encrypt short flags", "gpg -e -r user@example.com file.txt", "gpg -e -r <val> <path>"},
		{"decrypt with output", "gpg --output plaintext.txt --decrypt secret.txt.gpg", "gpg --output <path> --decrypt <path>"},
		{"decrypt short", "gpg -o plaintext.txt -d secret.txt.gpg", "gpg -o <path> -d <path>"},

		// Signing and verifying
		{"sign file", "gpg --sign document.pdf", "gpg --sign <path>"},
		{"sign detached", "gpg -ba filename", "gpg -ba <val>"},
		{"verify", "gpg --verify signature.asc document.pdf", "gpg --verify <path>+"},

		// Import / export
		{"import key", "gpg --import ~/public_key.txt", "gpg --import <path>"},
		{"export with armor", "gpg --armor --export ABCD1234", "gpg --armor --export <val>"},
		{"export secret key", "gpg --output ~/private_key.txt --armor --export-secret-key ABCD1234", "gpg --output <path> --armor --export-secret-key <val>"},

		// Key IDs (hex)
		{"fingerprint with key id", "gpg --fingerprint ABCD1234EF567890", "gpg --fingerprint <val>"},

		// Keyserver
		{"send keys with keyserver", "gpg --keyserver hkps://keys.openpgp.org --send-keys ABCD1234", "gpg --keyserver <val> --send-keys <val>"},
		{"search keys", "gpg --search-keys user@example.com", "gpg --search-keys <val>"},

		// Homedir
		{"custom homedir", "gpg --homedir /tmp/gnupg --list-keys", "gpg --homedir <path> --list-keys"},

		// gpg2 alias
		{"gpg2 alias", "gpg2 --list-keys", "gpg2 --list-keys"},
		{"gpg2 encrypt", "gpg2 -e -r user@example.com file.txt", "gpg2 -e -r <val> <path>"},

		// Default key
		{"default key", "gpg --default-key ABCD1234 -ba filename", "gpg --default-key <val> -ba <val>"},
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
	t.Run("different recipients collide", func(t *testing.T) {
		a := shellshape.Normalize("gpg -e -r alice@example.com file.txt")
		b := shellshape.Normalize("gpg -e -r bob@example.com file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different key ids collide", func(t *testing.T) {
		a := shellshape.Normalize("gpg --fingerprint ABCD1234")
		b := shellshape.Normalize("gpg --fingerprint DEADBEEF")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("gpg --encrypt file.txt")
		subshell := shellshape.Normalize("gpg --encrypt $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
