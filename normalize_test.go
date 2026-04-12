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
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple continuation", "echo \\\nhello", "echo <str>"},
		{"continuation with flags", "docker build \\\n  -t myapp \\\n  --no-cache \\\n  .", "docker build -t <val> --no-cache ."},
		{"continuation mid flag", "curl \\\n  -X POST \\\n  -d '{\"a\":1}' \\\n  https://api.example.com", "curl -X <method> -d <data> <https-uri>"},
		{"continuation preserves operators", "git add foo.ts && \\\ngit commit -m 'fix'", "git add <path> && git commit -m <str>"},
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

func TestNormalizeSubshellComplex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Recursive normalization of inner commands
		{"subshell with path", "cat $(find ./src -name foo.py)", "cat $(find <path> -name <pattern>)"},
		{"subshell inner normalized", "echo $(git rev-parse HEAD)", "echo $(git rev-parse HEAD)"},

		// Operators inside subshells stay contained
		{"subshell with pipe", "echo $(ls /tmp | head)", "echo $(ls <path> | head)"},
		{"subshell with and", "echo $(cd /tmp && ls)", "echo $(cd <path> && ls)"},
		{"subshell with semicolon", "echo $(echo hi; echo bye)", "echo $(echo <str> ; echo <str>)"},

		// Nested subshells
		{"nested subshell", "echo $(echo $(git status))", "echo $(echo $(git status))"},
		{"nested subshell with path", "echo $(cat $(find . -name foo.go))", "echo $(cat $(find . -name <pattern>))"},

		// Subshells with newlines inside
		{"subshell with newlines", "echo $(echo foo\necho bar)", "echo $(echo <str> ; echo <str>)"},

		// Multiple subshells in one command
		{"two subshells", "echo $(whoami) $(pwd)", "echo $(whoami) $(pwd)"},
		{"subshell and literal", "echo hello $(whoami)", "echo <str> $(whoami)"},

		// Subshell in pipeline
		{"subshell piped", "echo $(date) | cat", "echo $(date) | cat"},
		{"subshell after pipe", "cat file.txt | grep $(echo pattern)", "cat <path> | grep $(echo <str>)"},

		// Subshell in compound commands
		{"subshell with and operator", "echo $(whoami) && echo done", "echo $(whoami) && echo <str>"},
		{"subshell in second command", "cd /tmp && echo $(ls)", "cd <path> && echo $(ls)"},

		// Subshell in env var assignment
		{"subshell in env var", "A=$(date) cmd", "A=<val> cmd"},

		// Backticks
		{"backtick simple", "echo `date`", "echo $(<subshell>)"},
		{"backtick in pipeline", "echo `whoami` | cat", "echo $(<subshell>) | cat"},
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

func TestNormalizeHerestring(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple herestring", "cat <<< hello", "cat <<< <str>"},
		{"herestring quoted", `cat <<< "hello world"`, "cat <<< <str>"},
		{"herestring with grep", `grep pattern <<< "some text"`, "grep <pattern> <<< <str>"},
		{"herestring with jq", `jq '.name' <<< '{"name":"test"}'`, "jq <filter> <<< <str>"},
		{"herestring with wc", `wc -w <<< "count these words"`, "wc -w <<< <str>"},
		{"herestring with variable", "cat <<< $FOO", "cat <<< <str>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
	t.Run("different herestring values collide", func(t *testing.T) {
		a := Normalize(`cat <<< "hello"`)
		b := Normalize(`cat <<< "goodbye"`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
	t.Run("herestring differs from heredoc", func(t *testing.T) {
		a := Normalize("cat <<< hello")
		b := Normalize("cat <<EOF\nhello\nEOF")
		if a == b {
			t.Error("herestring and heredoc should produce different shapes")
		}
	})
}

func TestNormalizeCompoundCommands(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Mixed operators
		{"and then or", "pnpm build && pnpm test || echo failed", "pnpm build && pnpm test || echo <str>"},
		{"pipe then and", "git log | head && echo done", "git log | head && echo <str>"},
		{"semicolons and ands", "echo a; echo b && echo c", "echo <str> ; echo <str> && echo <str>"},

		// Newlines as separators mixed with operators
		{"newline then pipe", "git log\nls | head", "git log ; ls | head"},
		{"operator then newline", "git status &&\ngit log", "git status && git log"},
		{"newlines with all operators", "git add foo.ts\ngit commit -m 'fix' && git push\necho done", "git add <path> ; git commit -m <str> && git push ; echo <str>"},

		// Parenthesized subshells (not $())
		{"paren subshell", "(cd /tmp && ls)", "(cd <path> && ls)"},
		{"paren subshell piped", "(cd /tmp && ls) | head", "(cd <path> && ls) | head"},

		// Redirects with operators
		{"redirect then pipe", "grep pattern file.txt 2>/dev/null | head", "grep <pattern> <path> 2>/dev/null | head"},
		{"redirect then and", "make build 2>&1 && echo ok", "make build 2>&1 && echo <str>"},
		{"redirect in second command", "echo start && ls /tmp > out.txt", "echo <str> && ls <path> > <path>"},

		// Complex real-world patterns
		{"conditional with redirect", "test -f config.yml && cat config.yml || echo 'missing' > /dev/stderr", "test -f <path> && cat <path> || echo <str> > <path>"},
		{"pipeline with error handling", "curl -s https://api.example.com | jq '.data' || echo error", "curl -s <https-uri> | jq <filter> || echo <str>"},
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

func TestNormalizeHeredocComplex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Heredoc with operators before/after
		{"heredoc after and", "git add . && git commit -F - <<'EOF'\nfix bug\nEOF", "git add . && git commit -F - <heredoc>"},
		{"command after heredoc", "cat <<EOF\nhello\nEOF\necho done", "cat <heredoc> ; echo <str>"},

		// Heredoc followed by another command
		{"heredoc then command", "cat <<EOF\nhello world\nEOF\necho done", "cat <heredoc> ; echo <str>"},

		// Multiple heredocs (rare but valid)
		{"heredoc not confused with herestring", "cat <<EOF\ndata\nEOF", "cat <heredoc>"},
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

func TestNormalizeCoreEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// shelxSplit: backslash at end of input (line 66-69)
		{"trailing backslash", `echo hello\`, "echo <str>"},
		// shelxSplit: backslash in double quotes not escaping special (line 42-43)
		{"backslash non-special in dquotes", `echo "hello\nworld"`, "echo <str>"},

		// normalizeSingleCommand: env-only command (line 122)
		{"env only", "FOO=bar BAZ=qux", "FOO=<val> BAZ=<val>"},

		// normalizeSingleCommand: shell script runner with dotted file (line 142-144)
		{"bash dotted script", "bash script.sh arg1", "script.sh arg1"},

		// normalizeSingleCommand: redirect herestring (line 170-171)
		{"herestring redirect", "cat <<< hello", "cat <<< <str>"},

		// normalizeSingleCommand: redirect standalone (line 177-180)
		{"redirect standalone", "echo hello 2>&1", "echo <str> 2>&1"},

		// fallbackNormalize coverage
		{"unmatched single quote", "echo 'hello", "echo 'hello"},

		// normalizeSubstitutions: whole subshell is opaque
		{"nested subshell", "echo $(echo $(date))", "echo $(echo $(date))"},

		// splitRedirects: various redirect forms
		{"redirect with fd", "cmd 2> /tmp/err.log", "cmd 2> <path>"},

		// collapseRepeatedPlaceholders: non-adjacent same placeholders
		{"non-adjacent same", "cmd <path> -f <path>", "cmd <path> -f <path>"},

		// Shell-as-script-runner with flag (should not consume as script)
		{"bash with flag", "bash -c 'echo hello'", "bash -c <code>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeForLoop(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic for-in loop
		{"simple for loop", "for f in *.go; do echo $f; done", "for f in <val>+ ; do echo <str> ; done"},
		// For loop with multiple values
		{"for with values", "for cmd in doctor test show stats; do python3 app $cmd --help; done", "for cmd in <val>+ ; do python3 app <arg>+ ; done"},
		// For loop with && in body
		{"for with and in body", "for x in a b c; do echo $x && rm $x; done", "for x in <val>+ ; do echo <str> && rm $x ; done"},
		// For loop with pipe in body
		{"for with pipe in body", "for f in *.log; do cat $f | grep <pattern>; done", "for f in <val>+ ; do cat $f | grep <pattern> ; done"},
		// For loop chained with outer command
		{"for then and", "for f in a b; do echo $f; done && echo all done", "for f in <val>+ ; do echo <str> ; done && echo <str>"},
		// Nested for loops
		{"nested for", "for x in 1 2; do for y in a b; do echo $x $y; done; done", "for x in <val>+ ; do for y in <val>+ ; do echo <str> ; done ; done"},
		// C-style for loop (no in-list)
		{"c-style for", "for ((i=0; i<10; i++)); do echo $i; done", "for ((<expr>)) ; do echo <str> ; done"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeWhileUntilLoop(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"while loop", "while read line; do echo $line; done", "while read line ; do echo <str> ; done"},
		{"until loop", "until test -f /tmp/ready; do sleep 1; done", "until test -f <path> ; do sleep N ; done"},
		{"while with redirect", "while read line; do echo $line; done < /tmp/input", "while read line ; do echo <str> ; done < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeIfStatement(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple if", "if test -f /tmp/x; then echo yes; fi", "if test -f <path> ; then echo <str> ; fi"},
		{"if else", "if test -f /tmp/x; then echo yes; else echo no; fi", "if test -f <path> ; then echo <str> ; else echo <str> ; fi"},
		{"if elif", "if test -f a; then echo a; elif test -f b; then echo b; fi", "if test -f a ; then echo <str> ; elif test -f b ; then echo <str> ; fi"},
		{"if chained", "if test -f x; then echo y; fi && echo done", "if test -f x ; then echo <str> ; fi && echo <str>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q)\n  got  %q\n  want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExecutableOfEdgeCases(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"FOO=bar", ""},
		{"FOO=bar cmd", "cmd"},
		// ExecutableOf returns the raw token, not stripped
		{"/usr/bin/python3 script.py", "/usr/bin/python3"},
		{"cmd arg1 | cmd2 arg2", "cmd"},
		{"cmd arg1 && cmd2", "cmd"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ExecutableOf(tt.input)
			if got != tt.want {
				t.Errorf("ExecutableOf(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
