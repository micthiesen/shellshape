package shellshape

import "testing"

func TestClassifyToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Long flags
		{"long flag", "--oneline", "--oneline"},
		{"long flag with value", "--max-count=10", "--max-count=<val>"},
		{"long flag include glob", "--include=*.ts", "--include=<val>"},

		// Short flags
		{"short flag multi", "-rf", "-rf"},
		{"short flag single", "-v", "-v"},
		{"short flag digits", "-10", "N"},

		// Numbers
		{"integer", "42", "N"},
		{"float", "3.14", "N"},

		// URLs
		{"https url", "https://api.example.com/x?y=1", "<https-uri>"},
		{"s3 url", "s3://bucket/prefix/", "<s3-uri>"},

		// Git remote
		{"git remote", "git@github.com:owner/repo.git", "<git-uri>"},

		// UUID
		{"uuid", "550e8400-e29b-41d4-a716-446655440000", "<uuid>"},

		// Hash
		{"short hash", "48ee25e6b", "<hash>"},
		{"long hash", "48ee25e6b123456789012345678901234567890a", "<hash>"},

		// Paths
		{"relative path dot", "./build/cache", "<path>"},
		{"relative path dotdot", "../foo", "<path>"},
		{"absolute path", "/tmp/file.txt", "<path>"},
		{"home path", "~/code/foo", "<path>"},
		{"python file", "foo.py", "<path>"},
		{"json file", "settings.json", "<path>"},
		{"txt file via ext", "file.txt", "<path>"},

		// HEAD
		{"bare HEAD", "HEAD", "HEAD"},

		// Verbatim
		{"verbatim command", "status", "status"},

		// Dotted identifiers
		{"dotted python", "tests.test_foo", "<dotted-id>"},
		{"dotted java", "com.example.MyClass", "<dotted-id>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyToken(tt.input)
			if got != tt.want {
				t.Errorf("ClassifyToken(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsSubshellToken(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"$(echo hello)", true},
		{"$(date)", true},
		{"$VAR", false},
		{"echo", false},
		{"$(incomplete", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsSubshellToken(tt.input); got != tt.want {
				t.Errorf("IsSubshellToken(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsFlagToken(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"-v", true},
		{"-rf", true},
		{"--verbose", true},
		{"-10", false},
		{"---", false},
		{"hello", false},
		{"-", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsFlagToken(tt.input); got != tt.want {
				t.Errorf("IsFlagToken(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestIsEnvAssignment(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"FOO=bar", true},
		{"_X=1", true},
		{"PATH=/usr/bin", true},
		{"foo", false},
		{"=bar", false},
		{"1VAR=x", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := IsEnvAssignment(tt.input); got != tt.want {
				t.Errorf("IsEnvAssignment(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLooksLikeRev(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"HEAD", true},
		{"HEAD~3", true},
		{"HEAD^^", true},
		{"HEAD@{upstream}", true},
		{"abc1234", true},
		{"main", true},
		{"origin/main", true},
		{"feature-branch", true},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := looksLikeRev(tt.input); got != tt.want {
				t.Errorf("looksLikeRev(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
