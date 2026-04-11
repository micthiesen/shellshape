package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestMongo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare mongo", "mongo", "mongo"},
		{"bare mongosh", "mongosh", "mongosh"},
		{"connect to database", "mongo mydb", "mongo <db>"},
		{"mongosh connect to database", "mongosh mydb", "mongosh <db>"},
		{"connection string URI", "mongosh mongodb://user@host:27017/mydb", "mongosh <mongodb-uri>"},

		// Flags with arguments
		{"host and port", "mongosh --host localhost --port 27017 mydb", "mongosh --host <host> --port N <db>"},
		{"short host and port", "mongo -h 10.0.0.1 --port 27017", "mongo -h <host> --port N"},
		{"username and password", "mongo -u admin -p secret mydb", "mongo -u <val> -p <val> <db>"},
		{"long username", "mongosh --username admin --password secret mydb", "mongosh --username <val> --password <val> <db>"},
		{"eval expression", "mongosh --eval 'db.users.find()' mydb", "mongosh --eval <expr> <db>"},
		{"eval only", "mongo --eval 'JSON.stringify(db.foo.findOne())'", "mongo --eval <expr>"},
		{"auth database", "mongosh --username admin --authenticationDatabase admin mydb", "mongosh --username <val> --authenticationDatabase <val> <db>"},
		{"auth mechanism", "mongo --authenticationMechanism SCRAM-SHA-256 mydb", "mongo --authenticationMechanism <val> <db>"},

		// Script file positional (mongo legacy)
		{"db and script", "mongo localhost:27017/myDatabase script.js", "mongo <db> <path>"},
		{"mongosh file flag", "mongosh --file setup.js", "mongosh --file <path>"},
		{"mongosh short file", "mongosh -f init.js", "mongosh -f <path>"},

		// Boolean flags
		{"quiet mode", "mongosh --quiet mydb", "mongosh --quiet <db>"},
		{"shell mode", "mongo --shell script.js", "mongo --shell <db>"},
		{"nodb", "mongosh --nodb", "mongosh --nodb"},
		{"norc", "mongosh --norc mydb", "mongosh --norc <db>"},

		// TLS/SSL flags
		{"tls with cert", "mongosh --tls --tlsCertificateKeyFile /path/cert.pem mydb", "mongosh --tls --tlsCertificateKeyFile <path> <db>"},
		{"ssl cert", "mongo --ssl --sslCAFile /etc/ssl/ca.pem mydb", "mongo --ssl --sslCAFile <path> <db>"},
		{"tls key password", "mongosh --tlsCertificateKeyFilePassword secret123 mydb", "mongosh --tlsCertificateKeyFilePassword <val> <db>"},

		// Redirect
		{"with redirect", "mongo mydb > output.txt", "mongo <db> > <path>"},

		// Full realistic example
		{"full connection", "mongosh --host db.example.com --port 27017 --username admin --password s3cret --authenticationDatabase admin production", "mongosh --host <host> --port N --username <val> --password <val> --authenticationDatabase <val> <db>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q)\n  got:  %s\n  want: %s", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different databases collide", func(t *testing.T) {
		a := shellshape.Normalize("mongosh production")
		b := shellshape.Normalize("mongosh staging")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("mongo --host 10.0.0.1 --port 27017 mydb")
		b := shellshape.Normalize("mongo --host 192.168.1.5 --port 5555 otherdb")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	t.Run("different eval expressions collide", func(t *testing.T) {
		a := shellshape.Normalize("mongosh --eval 'db.users.find()'")
		b := shellshape.Normalize("mongosh --eval 'db.orders.count()'")
		if a != b {
			t.Errorf("expected same shape, got %q vs %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("mongosh literal-arg")
		subshell := shellshape.Normalize("mongosh $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
