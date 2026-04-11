package shellshape

import "testing"

func TestHTTPie(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple GET", "http https://example.com", "http <https-uri>"},
		{"https command", "https https://example.com/api", "https <https-uri>"},
		{"explicit GET", "http GET https://api.example.com/users", "http GET <https-uri>"},
		{"explicit POST", "http POST https://api.example.com/users name=John", "http POST <https-uri> <item>"},
		{"explicit PUT", "http PUT https://api.example.com/item/1 name=John", "http PUT <https-uri> <item>"},
		{"explicit DELETE", "http DELETE https://api.example.com/item/1", "http DELETE <https-uri>"},
		{"explicit PATCH", "http PATCH https://api.example.com/item/1 name=New", "http PATCH <https-uri> <item>"},

		// Request items (data)
		{"data field", "http POST https://example.com name=John email=john@example.com", "http POST <https-uri> <item>+"},
		{"json field", "http POST https://example.com age:=29 married:=false", "http POST <https-uri> <item>+"},
		{"header item", "http https://example.com X-API-Token:123", "http <https-uri> <item>"},
		{"query param", "http https://example.com search==term", "http <https-uri> <item>"},
		{"file upload", "http POST https://example.com cv@~/doc.pdf", "http POST <https-uri> <item>"},
		{"mixed items", "http PUT https://example.com name=John X-Token:abc age:=29", "http PUT <https-uri> <item>+"},
		{"file embed", "http POST https://example.com description=@about.txt", "http POST <https-uri> <item>"},

		// Flags with arguments
		{"auth flag", "http -a user:pass https://example.com", "http -a <data> <https-uri>"},
		{"long auth", "http --auth user:pass https://example.com", "http --auth <data> <https-uri>"},
		{"output flag", "http -o result.json https://example.com", "http -o <path> <https-uri>"},
		{"long output", "http --output result.json https://example.com", "http --output <path> <https-uri>"},
		{"session flag", "http --session=mysession https://example.com", "http --session=<path> <https-uri>"},
		{"session read-only", "http --session-read-only=test https://example.com", "http --session-read-only=<path> <https-uri>"},
		{"auth-type flag", "http --auth-type=digest -a user:pass https://example.com", "http --auth-type=<val> -a <data> <https-uri>"},
		{"proxy flag", "http --proxy=http:http://10.0.0.1:8080 https://example.com", "http --proxy=<val> <https-uri>"},
		{"verify flag", "http --verify=no https://example.com", "http --verify=<val> <https-uri>"},
		{"print flag", "http --print=hH https://example.com", "http --print=<val> <https-uri>"},
		{"style flag", "http --style monokai https://example.com", "http --style <val> <https-uri>"},

		// Numeric flags
		{"timeout flag", "http --timeout 30 https://example.com", "http --timeout N <https-uri>"},
		{"max-redirects", "http --max-redirects 5 https://example.com", "http --max-redirects N <https-uri>"},
		{"max-headers", "http --max-headers 100 https://example.com", "http --max-headers N <https-uri>"},

		// Boolean flags
		{"verbose", "http -v https://example.com", "http -v <https-uri>"},
		{"download", "http --download https://example.com/file.zip", "http --download <https-uri>"},
		{"form flag", "http -f POST https://example.com field=value", "http -f POST <https-uri> <item>"},
		{"json flag", "http --json POST https://example.com data=test", "http --json POST <https-uri> <item>"},
		{"headers only", "http -h https://example.com", "http -h <https-uri>"},
		{"body only", "http -b https://example.com", "http -b <https-uri>"},
		{"follow redirects", "http --follow https://example.com", "http --follow <https-uri>"},

		// Combined
		{"full example", "http -v -a admin:secret POST https://api.example.com/items name=Widget price:=9.99", "http -v -a <data> POST <https-uri> <item>+"},

		// Edge cases
		{"no args", "http", "http"},
		{"only flags", "http -v --json", "http -v --json"},
		{"cert flag", "http --cert /path/to/cert.pem https://example.com", "http --cert <path> <https-uri>"},
		{"cert-key flag", "http --cert-key /path/to/key.pem https://example.com", "http --cert-key <path> <https-uri>"},

		// Redirect
		{"output redirect", "http https://example.com > output.json", "http <https-uri> > <path>"},
		{"input redirect", "http POST https://example.com < data.json", "http POST <https-uri> < <path>"},
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
	t.Run("different URLs collide", func(t *testing.T) {
		a := Normalize("http GET https://api.example.com/users")
		b := Normalize("http GET https://api.other.com/items")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different request items collide", func(t *testing.T) {
		a := Normalize("http POST https://example.com name=Alice age:=30")
		b := Normalize("http POST https://example.com name=Bob age:=25")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different auth collide", func(t *testing.T) {
		a := Normalize("http -a alice:pass1 https://example.com")
		b := Normalize("http -a bob:pass2 https://example.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("http https://example.com")
		subshell := Normalize("http $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
