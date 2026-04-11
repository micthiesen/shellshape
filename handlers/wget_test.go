package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWget(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple URL", "wget https://example.com/file.tar.gz", "wget <https-uri>"},
		{"HTTP URL", "wget http://example.com/page.html", "wget <http-uri>"},
		{"FTP URL", "wget ftp://ftp.gnu.org/pub/file.tar.gz", "wget <ftp-uri>"},

		// Output flags
		{"output file", "wget -O output.tar.gz https://example.com/file.tar.gz", "wget -O <path> <https-uri>"},
		{"output to stdout", "wget -O - https://example.com/script.sh", "wget -O <path> <https-uri>"},
		{"log file", "wget -o wget.log https://example.com/file", "wget -o <path> <https-uri>"},
		{"append log", "wget -a wget.log https://example.com/file", "wget -a <path> <https-uri>"},

		// Numeric flags
		{"tries", "wget -t 3 https://example.com/file", "wget -t N <https-uri>"},
		{"timeout", "wget -T 30 https://example.com/file", "wget -T N <https-uri>"},
		{"wait", "wget -w 5 https://example.com/file", "wget -w N <https-uri>"},
		{"recursion depth", "wget -r -l 2 https://example.com", "wget -r -l N <https-uri>"},

		// Header and data flags
		{"header", "wget --header 'Authorization: Bearer token' https://api.example.com", "wget --header <header> <https-uri>"},
		{"post data", "wget --post-data 'key=value' https://api.example.com/endpoint", "wget --post-data <data> <https-uri>"},

		// String flags
		{"user agent short", "wget -U 'Mozilla/5.0' https://example.com", "wget -U <str> <https-uri>"},
		{"user agent long", "wget --user-agent 'CustomBot/1.0' https://example.com", "wget --user-agent <str> <https-uri>"},
		{"execute", "wget -e robots=off https://example.com", "wget -e <str> <https-uri>"},
		{"referer", "wget --referer https://google.com https://example.com/page", "wget --referer <str> <https-uri>"},

		// Val flags
		{"user", "wget --user admin https://example.com/private", "wget --user <val> <https-uri>"},
		{"password", "wget --password secret https://example.com/private", "wget --password <val> <https-uri>"},
		{"accept", "wget -A '*.jpg,*.png' https://example.com", "wget -A <val> <https-uri>"},
		{"reject", "wget -R '*.gif' https://example.com", "wget -R <val> <https-uri>"},
		{"domains", "wget -D example.com,example.org https://example.com", "wget -D <val> <https-uri>"},
		{"exclude dirs", "wget -X /cgi-bin,/tmp https://example.com", "wget -X <val> <https-uri>"},
		{"include dirs", "wget -I /docs https://example.com", "wget -I <val> <https-uri>"},
		{"method", "wget --method PUT https://api.example.com/resource", "wget --method <val> <https-uri>"},

		// Path flags
		{"input file", "wget -i urls.txt", "wget -i <path>"},
		{"directory prefix", "wget -P /tmp/downloads https://example.com/file", "wget -P <path> <https-uri>"},
		{"config", "wget --config /etc/wgetrc https://example.com", "wget --config <path> <https-uri>"},
		{"cookies load", "wget --load-cookies cookies.txt https://example.com", "wget --load-cookies <path> <https-uri>"},
		{"cookies save", "wget --save-cookies cookies.txt https://example.com", "wget --save-cookies <path> <https-uri>"},

		// Boolean flags
		{"quiet", "wget -q https://example.com/file", "wget -q <https-uri>"},
		{"continue", "wget -c https://example.com/large-file.iso", "wget -c <https-uri>"},
		{"recursive", "wget -r https://example.com", "wget -r <https-uri>"},
		{"mirror", "wget -m https://example.com", "wget -m <https-uri>"},
		{"spider", "wget --spider https://example.com/check", "wget --spider <https-uri>"},
		{"no clobber", "wget -nc https://example.com/file", "wget -nc <https-uri>"},
		{"timestamping", "wget -N https://example.com/file", "wget -N <https-uri>"},

		// Combined flags
		{"mirror combo", "wget -m -p -k https://example.com", "wget -m -p -k <https-uri>"},
		{"recursive with depth and wait", "wget -r -l 2 -w 5 https://example.com", "wget -r -l N -w N <https-uri>"},
		{"quiet output to file", "wget -q -O /tmp/file https://example.com/data", "wget -q -O <path> <https-uri>"},

		// Multiple URLs
		{"multiple URLs", "wget https://example.com/a https://example.com/b", "wget <https-uri>+"},

		// Redirects
		{"with redirect", "wget -q https://example.com/file > /dev/null", "wget -q <https-uri> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different URLs collide", func(t *testing.T) {
		a := shellshape.Normalize("wget https://example.com/file1.tar.gz")
		b := shellshape.Normalize("wget https://other-site.org/different-file.zip")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different output files collide", func(t *testing.T) {
		a := shellshape.Normalize("wget -O /tmp/output1.txt https://example.com/a")
		b := shellshape.Normalize("wget -O /home/user/result.json https://example.com/b")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("wget https://example.com/file")
		subshell := shellshape.Normalize("wget $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
