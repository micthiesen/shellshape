package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestOpenssl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"req self-signed", "openssl req -new -x509 -key server.key -out server.crt -days 365", "openssl req -new -x509 -key <path> -out <path> -days N"},
		{"x509 inspect", "openssl x509 -in cert.pem -text -noout", "openssl x509 -in <path> -text -noout"},
		{"s_client connect", "openssl s_client -connect example.com:443", "openssl s_client -connect <val>"},
		{"genrsa", "openssl genrsa -out private.key 2048", "openssl genrsa -out <path> N"},
		{"enc encrypt", "openssl enc -aes-256-cbc -in file.txt -out file.enc -pass pass:mypassword", "openssl enc -aes-256-cbc -in <path> -out <path> -pass <val>"},
		{"dgst sha256", "openssl dgst -sha256 file.txt", "openssl dgst -sha256 <path>"},
		{"pkcs12 export", "openssl pkcs12 -export -in cert.pem -inkey key.pem -out cert.p12", "openssl pkcs12 -export -in <path> -inkey <path> -out <path>"},
		{"rsa pubout", "openssl rsa -in private.key -pubout -out public.key", "openssl rsa -in <path> -pubout -out <path>"},
		{"verify with CA", "openssl verify -CAfile ca.pem cert.pem", "openssl verify -CAfile <path>+"},
		{"rand base64", "openssl rand -base64 32", "openssl rand -base64 N"},
		{"req with subj", "openssl req -new -key server.key -subj /CN=example.com -out server.csr", "openssl req -new -key <path> -subj <val> -out <path>"},
		{"version", "openssl version", "openssl version"},
		{"x509 with days", "openssl x509 -req -in server.csr -CA ca.pem -CAkey ca-key.pem -out server.crt -days 365", "openssl x509 -req -in <path> -CA <path> -CAkey <path> -out <path> -days N"},
		{"s_client with cert", "openssl s_client -connect host:443 -cert client.pem -key client-key.pem", "openssl s_client -connect <val> -cert <path> -key <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("openssl x509 -in /tmp/cert1.pem -text -noout")
		b := shellshape.Normalize("openssl x509 -in /etc/ssl/cert2.pem -text -noout")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different keys collide", func(t *testing.T) {
		a := shellshape.Normalize("openssl genrsa -out mykey.key 4096")
		b := shellshape.Normalize("openssl genrsa -out otherkey.key 2048")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("openssl dgst -sha256 file.txt")
		subshell := shellshape.Normalize("openssl dgst -sha256 $(find . -name cert.pem)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
