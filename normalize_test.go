package shellshape

import "testing"

func TestNormalizeSimple(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"whitespace", "   ", ""},
		{"bare command", "ls", "ls"},
		{"ls with path", "ls /tmp", "ls <path>"},
		{"git log oneline", "git log --oneline -10", "git log --oneline N"},
		{"git log count flag", "git log --max-count=20", "git log --max-count=<val>"},
		{"git status", "git status", "git status"},
		{"rm rf relative", "rm -rf ./build", "rm -rf <path>"},
		{"rm rf absolute", "rm -rf /tmp/foo", "rm -rf <path>"},
		{"pnpm test file", "pnpm test packages/domain/foo.spec.ts", "pnpm test <path>"},
		{"aws s3 ls", "aws s3 ls s3://my-bucket/prefix/", "aws s3 ls <s3-uri>"},
		{"curl url", "curl https://api.foo.com/x", "curl <https-uri>"},
		{"whitespace collapsed", "git    log     --oneline", "git log --oneline"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeEnvVars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"single env var", "FOO=bar pnpm build", "FOO=<val> pnpm build"},
		{"multiple env vars", "FOO=1 BAR=baz node index.js", "FOO=<val> BAR=<val> node <path>"},
		{"env var only", "FOO=bar", "FOO=<val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizePipes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple pipe", "git log | head", "git log | head"},
		{"pipe with args", "git log --oneline | head -5", "git log --oneline | head N"},
		{"double and", "pnpm test && pnpm build", "pnpm test && pnpm build"},
		{"semicolon", "cd /tmp; ls", "cd <path> ; ls"},
		{"or", "pnpm test || echo failed", "pnpm test || echo <str>"},
		{"three stage pipe", "cat log.txt | grep error | wc -l", "cat <path> | grep <pattern> | wc -l"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeSubstitution(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"dollar paren", "echo $(git rev-parse HEAD)", "echo $(git rev-parse HEAD)"},
		{"dollar paren with path", "cat $(find ./src -name foo.py)", "cat $(find <path> -name <pattern>)"},
		{"backtick", "echo `date`", "echo $(<subshell>)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeRedirects(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"stdout redirect", "echo hi > /tmp/out.txt", "echo <str> > <path>"},
		{"stderr redirect", "node build.js 2> build.log", "node <path> 2> <path>"},
		{"combined redirect", "node build.js 2>&1", "node <path> 2>&1"},
		{"append redirect", "echo hi >> /tmp/out", "echo <str> >> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeQuoted(t *testing.T) {
	t.Run("double quoted string", func(t *testing.T) {
		result := Normalize(`git commit -m "hello world"`)
		if !contains(result, "git commit -m") {
			t.Errorf("got %q, want it to contain 'git commit -m'", result)
		}
	})
	t.Run("single quoted grep pattern", func(t *testing.T) {
		got := Normalize("grep 'error' /var/log/foo.log")
		want := "grep <pattern> <path>"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestNormalizeMultiline(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"newline split", "git status\ngit log", "git status ; git log"},
		{"three lines", "git status\ngit log\necho hi", "git status ; git log ; echo <str>"},
		{"crlf", "git status\r\ngit log", "git status ; git log"},
		{"blank lines collapsed", "git status\n\n\ngit log", "git status ; git log"},
		{"newline with pipes", "git log | head\nls -la", "git log | head ; ls -la"},
		{"trailing newline", "git status\n", "git status"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeTricky(t *testing.T) {
	t.Run("unbalanced quotes fallthrough", func(t *testing.T) {
		result := Normalize("echo 'unclosed")
		if result == "" {
			t.Error("expected non-empty result for unbalanced quotes")
		}
	})
	t.Run("stable across argument differences", func(t *testing.T) {
		a := Normalize("rm -rf ./node_modules")
		b := Normalize("rm -rf /tmp/xyz")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
	t.Run("stable across number differences", func(t *testing.T) {
		a := Normalize("git log -5")
		b := Normalize("git log -50")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
	t.Run("different subcommands different shape", func(t *testing.T) {
		if Normalize("git push") == Normalize("git status") {
			t.Error("git push and git status should have different shapes")
		}
	})
	t.Run("different executables different shape", func(t *testing.T) {
		if Normalize("rm foo") == Normalize("ls foo") {
			t.Error("rm and ls should have different shapes")
		}
	})
}

func TestNormalizeComments(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"full line comment dropped", "# run the tests\npnpm test", "pnpm test"},
		{"trailing comment dropped", "pnpm test # run the tests", "pnpm test"},
		{"comment eats rest of line", "git status ; # comment eats ; git log", "git status"},
		{"comment after newline", "git status\n# a comment\ngit log", "git status ; git log"},
		{"only comment returns empty", "# just a comment", ""},
		{"comment at start after operator", "git status && # comment\ngit log", "git status && git log"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("hash inside single quotes not comment", func(t *testing.T) {
		result := Normalize("echo '#foo'")
		if !contains(result, "echo") {
			t.Errorf("got %q, expected to contain 'echo'", result)
		}
		if result == "echo" {
			t.Error("hash inside quotes should not strip content")
		}
	})
	t.Run("hash inside double quotes not comment", func(t *testing.T) {
		result := Normalize(`echo "#foo bar"`)
		if !contains(result, "echo") {
			t.Errorf("got %q, expected to contain 'echo'", result)
		}
		if result == "echo" {
			t.Error("hash inside quotes should not strip content")
		}
	})
	t.Run("hash mid word not comment", func(t *testing.T) {
		result := Normalize("git log --pretty=#short")
		if !contains(result, "git log") {
			t.Errorf("got %q, expected to contain 'git log'", result)
		}
	})
}

func TestNormalizeRevRanges(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"sha range to HEAD", "git log abc1234..HEAD", "git log <range>"},
		{"full sha range to main", "git log 64657e2cbbbc11916fc9c037e02e6310020c36fe..main", "git log <range>"},
		{"sha to sha", "git diff abc1234..def5678", "git diff <range>"},
		{"three dot range", "git log main...feature", "git log <range>"},
		{"branch to branch", "git log feature1..main", "git log <range>"},
		{"origin slash branch range", "git log origin/main..feature", "git log <range>"},
		{"HEAD tilde n", "git reset --hard HEAD~5", "git reset --hard <rev>"},
		{"HEAD caret", "git show HEAD^", "git show <rev>"},
		{"HEAD double caret", "git show HEAD^^", "git show <rev>"},
		{"bare HEAD kept", "git show HEAD", "git show HEAD"},
		{"range with path filter", "git log abc1234..HEAD -- src/foo.py", "git log <range> -- <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different SHAs collide", func(t *testing.T) {
		a := Normalize("git log abc1234..main")
		b := Normalize("git log def5678..main")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeInterpreterCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"python -c", `python3 -c "import sys; print(sys.version)"`, "python3 -c <code>"},
		{"bash -c", "bash -c 'ls && cd /tmp'", "bash -c <code>"},
		{"sh -c", "sh -c 'echo hi'", "sh -c <code>"},
		{"node -e", `node -e "console.log(1)"`, "node -e <code>"},
		{"python script", "python3 script.py", "python3 <path>"},
		{"bash -x is flag not code", "bash -x script.sh", "bash -x <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different code bodies collide", func(t *testing.T) {
		a := Normalize(`python3 -c "import sys"`)
		b := Normalize(`python3 -c "print(1+1)"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeAbsoluteExecutable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"absolute path exe", "/Users/michael/.dotfiles/claude/hooks/approvals/approvals show --always", "approvals show --always"},
		{"absolute path exe with args", "/usr/local/bin/python3 -c 'print(1)'", "python3 -c <code>"},
		{"tilde path exe", "~/.dotfiles/claude/hooks/approvals/approvals events --limit 5", "approvals events --limit N"},
		{"tilde path exe with subcommand", "~/.local/bin/gh pr list", "gh pr list"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("tilde and absolute paths collide", func(t *testing.T) {
		a := Normalize("~/.dotfiles/bin/approvals doctor")
		b := Normalize("/usr/local/bin/approvals doctor")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
	t.Run("different absolute paths collide", func(t *testing.T) {
		a := Normalize("/opt/homebrew/bin/approvals doctor")
		b := Normalize("/Users/michael/.local/bin/approvals doctor")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeBashPrefix(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bash script path", "bash .claude/skills/cut-release/scripts/gather.sh", "gather.sh"},
		{"bash relative script", "bash ./scripts/deploy.sh", "deploy.sh"},
		{"bash absolute script", "bash /Users/michael/.dotfiles/scripts/setup.sh", "setup.sh"},
		{"sh script", "sh ./run-tests.sh", "run-tests.sh"},
		{"zsh script", "zsh ~/.dotfiles/install.zsh", "install.zsh"},
		{"bash script with args", "bash ./deploy.sh --env production --dry-run", "deploy.sh --env production --dry-run"},
		{"bash -c not affected", "bash -c 'echo hello'", "bash -c <code>"},
		{"bare bash", "bash", "bash"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different paths same script collide", func(t *testing.T) {
		a := Normalize("bash /home/user/scripts/gather.sh")
		b := Normalize("bash .claude/skills/release/gather.sh")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
	t.Run("different scripts dont collide", func(t *testing.T) {
		a := Normalize("bash ./gather.sh")
		b := Normalize("bash ./deploy.sh")
		if a == b {
			t.Error("different scripts should have different shapes")
		}
	})
}

func TestNormalizeDottedPaths(t *testing.T) {
	t.Run("python unittest module", func(t *testing.T) {
		got := Normalize("python3 -m unittest tests.test_normalize.TestNormalizeFind")
		want := "python3 -m unittest <dotted-id>"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	t.Run("multiple dotted args collide", func(t *testing.T) {
		a := Normalize("python3 -m unittest tests.test_normalize.TestNormalizeFind tests.test_normalize.TestNormalizeGrep")
		b := Normalize("python3 -m unittest tests.test_cache.TestCacheHit tests.test_decision.TestVerdictSerialization")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeRepeatedPlaceholders(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"git add two paths", "git add foo.ts bar.ts", "git add <path>+"},
		{"git add many paths", "git add a.ts b.ts c.ts d.ts e.ts f.ts", "git add <path>+"},
		{"single path not collapsed", "git add foo.ts", "git add <path>"},
		{"grep multiple files", "grep -rn pattern /a /b /c /d", "grep -rn <pattern> <path>+"},
		{"non adjacent not collapsed", "cp foo.ts bar.ts && mv baz.ts qux.ts", "cp <path>+ && mv <path>+"},
		{"cat with redirect", "cat foo.ts > bar.ts", "cat <path> > <path>"},
		{"collapse in pipeline", "grep -rn pattern /a /b 2>/dev/null | head -5", "grep -rn <pattern> <path>+ 2>/dev/null | head N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different counts collide", func(t *testing.T) {
		a := Normalize("git add a.ts b.ts")
		b := Normalize("git add a.ts b.ts c.ts d.ts e.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeHeredoc(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"git commit heredoc", "git commit -F - <<'EOF'\nFix the bug\nEOF", "git commit -F - <heredoc>"},
		{"git commit heredoc unquoted", "git commit -F - <<EOF\nFix the bug\nEOF", "git commit -F - <heredoc>"},
		{"heredoc multiline", "git commit -F - <<'EOF'\nline one\nline two\nline three\nEOF", "git commit -F - <heredoc>"},
		{"heredoc with pipe", "git add foo.ts && git commit -F - <<'EOF'\nFix\nEOF", "git add <path> && git commit -F - <heredoc>"},
		{"heredoc indented", "cat <<-EOF\n\thello\n\tEOF", "cat <heredoc>"},
		{"pbcopy heredoc", "pbcopy <<EOF\nsome text\nEOF", "pbcopy <heredoc>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different commit messages collide", func(t *testing.T) {
		a := Normalize("git commit -F - <<'EOF'\nClean up VJ late mending\nEOF")
		b := Normalize("git commit -F - <<'EOF'\nReduce noise from IndexedDB\nEOF")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}

func TestNormalizeLineContinuation(t *testing.T) {
	t.Run("simple continuation", func(t *testing.T) {
		got := Normalize("echo \\\nhello")
		want := "echo <str>"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestExecutableOf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "git status", "git"},
		{"with env vars", "FOO=<val> pnpm build", "pnpm"},
		{"with pipe", "git log | head", "git"},
		{"with rm", "rm -rf <path>", "rm"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExecutableOf(tt.input)
			if got != tt.want {
				t.Errorf("ExecutableOf(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSplitTopLevel(t *testing.T) {
	t.Run("no operator", func(t *testing.T) {
		parts := splitTopLevel("git status", []string{"|", "&&"})
		if len(parts) != 1 || parts[0].text != "git status" || parts[0].sep != "" {
			t.Errorf("got %+v", parts)
		}
	})
	t.Run("single pipe", func(t *testing.T) {
		parts := splitTopLevel("a | b", []string{"|"})
		if len(parts) != 2 || parts[0].text != "a " || parts[0].sep != "|" || parts[1].text != " b" || parts[1].sep != "" {
			t.Errorf("got %+v", parts)
		}
	})
	t.Run("pipe inside quotes ignored", func(t *testing.T) {
		parts := splitTopLevel(`echo "a | b"`, []string{"|"})
		if len(parts) != 1 {
			t.Errorf("expected 1 segment, got %d: %+v", len(parts), parts)
		}
	})
	t.Run("pipe inside subshell ignored", func(t *testing.T) {
		parts := splitTopLevel("echo $(a | b)", []string{"|"})
		if len(parts) != 1 {
			t.Errorf("expected 1 segment, got %d: %+v", len(parts), parts)
		}
	})
	t.Run("longest match wins", func(t *testing.T) {
		parts := splitTopLevel("a || b", []string{"||", "|"})
		if len(parts) != 2 || parts[0].sep != "||" {
			t.Errorf("expected || separator, got %+v", parts)
		}
	})
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
