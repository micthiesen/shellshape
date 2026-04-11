package shellshape

import "testing"

func TestCurl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple GET", "curl https://example.com", "curl <https-uri>"},
		{"http URL", "curl http://example.com/path", "curl <http-uri>"},
		{"silent GET", "curl -s https://api.example.com/data", "curl -s <https-uri>"},

		// Flags with arguments
		{"output file", "curl -o output.json https://example.com", "curl -o <path> <https-uri>"},
		{"POST with data", "curl -X POST -d '{\"key\":\"val\"}' https://api.example.com", "curl -X <method> -d <data> <https-uri>"},
		{"custom header", "curl -H 'Authorization: Bearer token123' https://api.example.com", "curl -H <header> <https-uri>"},
		{"multiple headers", "curl -H 'Content-Type: application/json' -H 'Accept: application/json' https://api.example.com", "curl -H <header> -H <header> <https-uri>"},
		{"user auth", "curl -u admin:secret https://api.example.com", "curl -u <data> <https-uri>"},
		{"user agent", "curl -A 'Mozilla/5.0' https://example.com", "curl -A <str> <https-uri>"},
		{"form upload", "curl -F 'file=@photo.jpg' https://api.example.com/upload", "curl -F <data> <https-uri>"},
		{"json data", "curl --json '{\"name\":\"test\"}' https://api.example.com", "curl --json <data> <https-uri>"},
		{"write-out format", "curl -w '%{http_code}' https://example.com", "curl -w <val> <https-uri>"},

		// Numeric flags
		{"max time", "curl -m 30 https://example.com", "curl -m N <https-uri>"},
		{"retry", "curl --retry 3 https://example.com", "curl --retry N <https-uri>"},
		{"max redirs", "curl -L --max-redirs 10 https://example.com", "curl -L --max-redirs N <https-uri>"},
		{"connect timeout", "curl --connect-timeout 5 https://example.com", "curl --connect-timeout N <https-uri>"},

		// Combined flags
		{"silent show-error", "curl -sS https://example.com", "curl -sS <https-uri>"},
		{"follow redirects verbose", "curl -Lv https://example.com", "curl -Lv <https-uri>"},
		{"full POST example", "curl -s -X POST -H 'Content-Type: application/json' -d '{\"q\":\"test\"}' -o /tmp/resp.json https://api.example.com/search", "curl -s -X <method> -H <header> -d <data> -o <path> <https-uri>"},

		// Long form flags with arguments
		{"long data", "curl --data 'param=value' https://example.com", "curl --data <data> <https-uri>"},
		{"long header", "curl --header 'X-Custom: foo' https://example.com", "curl --header <header> <https-uri>"},
		{"long request", "curl --request DELETE https://api.example.com/item/42", "curl --request <method> <https-uri>"},
		{"long output", "curl --output /tmp/file.html https://example.com", "curl --output <path> <https-uri>"},

		// Path/file flags
		{"upload file", "curl -T /path/to/file.txt ftp://example.com", "curl -T <path> <ftp-uri>"},
		{"cookie file", "curl -b /tmp/cookies.txt https://example.com", "curl -b <path> <https-uri>"},
		{"cookie jar", "curl -c /tmp/cookies.txt https://example.com", "curl -c <path> <https-uri>"},
		{"cacert", "curl --cacert /etc/ssl/cert.pem https://example.com", "curl --cacert <path> <https-uri>"},
		{"dump header", "curl -D /tmp/headers.txt https://example.com", "curl -D <path> <https-uri>"},

		// Multiple URLs
		{"multiple URLs", "curl https://example.com https://other.com", "curl <https-uri>+"},

		// Edge cases
		{"no args", "curl", "curl"},
		{"only flags", "curl -v -s", "curl -v -s"},
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
		a := Normalize("curl -s https://api.example.com/users")
		b := Normalize("curl -s https://api.other.com/items")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different data collide", func(t *testing.T) {
		a := Normalize("curl -d 'user=alice' https://example.com")
		b := Normalize("curl -d 'user=bob' https://example.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different headers collide", func(t *testing.T) {
		a := Normalize("curl -H 'Authorization: Bearer abc' https://example.com")
		b := Normalize("curl -H 'Authorization: Bearer xyz' https://example.com")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("curl https://example.com")
		subshell := Normalize("curl $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
