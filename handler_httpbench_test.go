package shellshape

import "testing"

func TestAb(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple URL", "ab https://example.com/", "ab <https-uri>"},
		{"http URL", "ab http://example.com/path", "ab <http-uri>"},

		// Numeric flags
		{"requests and concurrency", "ab -n 100 -c 10 https://example.com/", "ab -n N -c N <https-uri>"},
		{"timelimit", "ab -t 30 -c 50 https://example.com/", "ab -t N -c N <https-uri>"},
		{"verbosity", "ab -v 2 https://example.com/", "ab -v N <https-uri>"},
		{"socket timeout", "ab -s 10 https://example.com/", "ab -s N <https-uri>"},
		{"buffer size", "ab -b 4096 https://example.com/", "ab -b N <https-uri>"},

		// Path flags
		{"csv output", "ab -e /tmp/results.csv https://example.com/", "ab -e <path> <https-uri>"},
		{"gnuplot output", "ab -g /tmp/results.tsv https://example.com/", "ab -g <path> <https-uri>"},
		{"POST file", "ab -n 100 -T application/json -p /tmp/data.json https://api.example.com/", "ab -n N -T <val> -p <path> <https-uri>"},
		{"PUT file", "ab -n 50 -u /tmp/data.txt https://api.example.com/", "ab -n N -u <path> <https-uri>"},
		{"client cert", "ab -E /etc/ssl/cert.pem https://example.com/", "ab -E <path> <https-uri>"},

		// Data flags (auth/cookie)
		{"basic auth", "ab -A admin:secret https://example.com/", "ab -A <data> <https-uri>"},
		{"cookie", "ab -C session=abc123 https://example.com/", "ab -C <data> <https-uri>"},
		{"proxy auth", "ab -P proxyuser:pass https://example.com/", "ab -P <data> <https-uri>"},

		// Header flag
		{"custom header", "ab -H 'Accept-Encoding: gzip' https://example.com/", "ab -H <header> <https-uri>"},

		// Val flags
		{"content type", "ab -T application/json https://example.com/", "ab -T <val> <https-uri>"},
		{"http method", "ab -m PUT https://example.com/", "ab -m <val> <https-uri>"},
		{"proxy", "ab -X proxy.example.com:8080 https://example.com/", "ab -X <val> <https-uri>"},
		{"cipher suite", "ab -Z DHE-RSA-AES256-SHA https://example.com/", "ab -Z <val> <https-uri>"},
		{"ssl protocol", "ab -f TLS1.2 https://example.com/", "ab -f <val> <https-uri>"},
		{"bind address", "ab -B 192.168.1.100 https://example.com/", "ab -B <val> <https-uri>"},

		// Boolean flags
		{"keepalive", "ab -k https://example.com/", "ab -k <https-uri>"},
		{"keepalive with requests", "ab -k -n 1000 -c 50 https://example.com/", "ab -k -n N -c N <https-uri>"},
		{"quiet", "ab -q -n 10000 https://example.com/", "ab -q -n N <https-uri>"},

		// Combined
		{"full example", "ab -n 100 -c 10 -H 'Host: app.example.com' -T application/json -p /tmp/body.json https://app.example.com/api", "ab -n N -c N -H <header> -T <val> -p <path> <https-uri>"},

		// Edge cases
		{"no args", "ab", "ab"},
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
		a := Normalize("ab -n 100 https://api.example.com/users")
		b := Normalize("ab -n 100 https://api.other.com/items")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different auth collide", func(t *testing.T) {
		a := Normalize("ab -A alice:pass1 https://example.com/")
		b := Normalize("ab -A bob:pass2 https://example.com/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ab https://example.com/")
		subshell := Normalize("ab $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestWrk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple URL", "wrk https://example.com/", "wrk <https-uri>"},

		// Numeric flags
		{"fused flags", "wrk -t12 -c400 -d30s http://127.0.0.1:8080/index.html", "wrk -t12 -c400 -d30s <http-uri>"},
		{"long form threads", "wrk --threads 4 --connections 100 https://example.com/", "wrk --threads N --connections N <https-uri>"},

		// Duration flags
		{"duration", "wrk -d 30s https://example.com/", "wrk -d <val> <https-uri>"},
		{"timeout", "wrk --timeout 5s https://example.com/", "wrk --timeout <val> <https-uri>"},

		// Header flag
		{"custom header", "wrk -H 'Host: example.com' https://example.com/", "wrk -H <header> <https-uri>"},

		// Script flag
		{"lua script", "wrk -s /path/to/script.lua https://example.com/", "wrk -s <path> <https-uri>"},
		{"long script", "wrk --script /path/to/script.lua https://example.com/", "wrk --script <path> <https-uri>"},

		// Boolean flags
		{"latency", "wrk --latency https://example.com/", "wrk --latency <https-uri>"},

		// Combined
		{"full example", "wrk -t4 -c100 -d30s --latency -H 'Accept: application/json' -s /tmp/post.lua https://api.example.com/", "wrk -t4 -c100 -d30s --latency -H <header> -s <path> <https-uri>"},

		// Edge cases
		{"no args", "wrk", "wrk"},
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
		a := Normalize("wrk -t4 -c100 -d30s https://api.example.com/users")
		b := Normalize("wrk -t4 -c100 -d30s https://api.other.com/items")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("wrk https://example.com/")
		subshell := Normalize("wrk $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestHey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple URL", "hey https://example.com/", "hey <https-uri>"},
		{"http URL", "hey http://localhost:8080/api", "hey <http-uri>"},

		// Numeric flags
		{"requests", "hey -n 200 https://example.com/", "hey -n N <https-uri>"},
		{"concurrency", "hey -c 50 https://example.com/", "hey -c N <https-uri>"},
		{"rate limit", "hey -q 10 https://example.com/", "hey -q N <https-uri>"},
		{"timeout", "hey -t 30 https://example.com/", "hey -t N <https-uri>"},
		{"cpus", "hey -cpus 4 https://example.com/", "hey -cpus N <https-uri>"},

		// Duration flag
		{"duration", "hey -z 10s https://example.com/", "hey -z <val> <https-uri>"},

		// Header flag
		{"custom header", "hey -H 'Authorization: Bearer token123' https://api.example.com/", "hey -H <header> <https-uri>"},

		// Data flags
		{"request body", "hey -d '{\"key\":\"value\"}' https://api.example.com/", "hey -d <data> <https-uri>"},
		{"basic auth", "hey -a admin:secret https://example.com/", "hey -a <data> <https-uri>"},

		// Path flag
		{"body from file", "hey -D /tmp/body.json https://api.example.com/", "hey -D <path> <https-uri>"},

		// Val flags
		{"http method", "hey -m POST https://api.example.com/", "hey -m <val> <https-uri>"},
		{"content type", "hey -T application/json https://api.example.com/", "hey -T <val> <https-uri>"},
		{"accept header", "hey -A application/json https://api.example.com/", "hey -A <val> <https-uri>"},
		{"output format", "hey -o csv https://example.com/", "hey -o <val> <https-uri>"},
		{"proxy", "hey -x http://proxy:8080 https://example.com/", "hey -x <val> <https-uri>"},
		{"host header", "hey -host app.example.com https://example.com/", "hey -host <val> <https-uri>"},

		// Boolean flags
		{"http2", "hey -h2 https://example.com/", "hey -h2 <https-uri>"},
		{"disable compression", "hey -disable-compression https://example.com/", "hey -disable-compression <https-uri>"},
		{"disable keepalive", "hey -disable-keepalive https://example.com/", "hey -disable-keepalive <https-uri>"},
		{"disable redirects", "hey -disable-redirects https://example.com/", "hey -disable-redirects <https-uri>"},

		// Combined
		{"full POST", "hey -n 500 -c 20 -m POST -T application/json -d '{\"q\":\"test\"}' -H 'X-Request-Id: 123' https://api.example.com/search", "hey -n N -c N -m <val> -T <val> -d <data> -H <header> <https-uri>"},

		// Edge cases
		{"no args", "hey", "hey"},
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
		a := Normalize("hey -n 200 https://api.example.com/users")
		b := Normalize("hey -n 200 https://api.other.com/items")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different body data collide", func(t *testing.T) {
		a := Normalize("hey -d '{\"name\":\"alice\"}' https://example.com/")
		b := Normalize("hey -d '{\"name\":\"bob\"}' https://example.com/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("hey https://example.com/")
		subshell := Normalize("hey $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
