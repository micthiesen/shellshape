package shellshape

import "testing"

func TestGh(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Second subcommand kept verbatim
		{"pr create basic", "gh pr create --title 'Fix bug' --body 'Description'", "gh pr create --title <str> --body <str>"},
		{"pr view number", "gh pr view 123", "gh pr view N"},
		{"pr view web", "gh pr view 123 --web", "gh pr view N --web"},
		{"pr checkout", "gh pr checkout 456", "gh pr checkout N"},
		{"pr list with flags", "gh pr list --state open --label bug --assignee octocat", "gh pr list --state <val> --label <val> --assignee <val>"},
		{"pr merge number", "gh pr merge 42 --squash", "gh pr merge N --squash"},
		{"pr diff", "gh pr diff 99", "gh pr diff N"},
		{"pr list json jq", "gh pr list --json number,title --jq '.[].number'", "gh pr list --json <val> --jq <val>"},

		// Issue commands
		{"issue create", "gh issue create -t 'Bug report' -b 'Steps to reproduce'", "gh issue create -t <str> -b <str>"},
		{"issue view", "gh issue view 789", "gh issue view N"},
		{"issue list labels", "gh issue list --label bug --label urgent", "gh issue list --label <val> --label <val>"},
		{"issue close", "gh issue close 55", "gh issue close N"},

		// Repo commands
		{"repo clone", "gh repo clone owner/repo", "gh repo clone <path>"},
		{"repo view web", "gh repo view --web", "gh repo view --web"},
		{"repo create", "gh repo create my-repo --public", "gh repo create my-repo --public"},
		{"repo fork", "gh repo fork owner/repo --clone", "gh repo fork <path> --clone"},

		// API command
		{"api endpoint", "gh api repos/owner/repo/pulls", "gh api <path>"},
		{"api with jq", "gh api repos/owner/repo/pulls --jq '.[].title'", "gh api <path> --jq <val>"},
		{"api method", "gh api repos/owner/repo/issues -X POST -f title=hello", "gh api <path> -X <val> -f <val>"},

		// Short flags
		{"short repo flag", "gh pr view 10 -R owner/repo", "gh pr view N -R <val>"},
		{"short web flag", "gh issue view 5 -w", "gh issue view N -w"},

		// PR view by URL
		{"pr view url", "gh pr view https://github.com/owner/repo/pull/123", "gh pr view <https-uri>"},

		// Draft flag
		{"pr create draft", "gh pr create --title 'WIP' --body 'draft' --draft", "gh pr create --title <str> --body <str> --draft"},

		// Limit flag
		{"pr list limit", "gh pr list --limit 50", "gh pr list --limit N"},
		{"issue list limit", "gh issue list -L 20", "gh issue list -L N"},

		// pr create with base/head
		{"pr create base head", "gh pr create --base main --head feature-branch", "gh pr create --base <val> --head <val>"},

		// pr create with reviewers
		{"pr create reviewers", "gh pr create --title 'PR' --reviewer alice --reviewer bob", "gh pr create --title <str> --reviewer <val> --reviewer <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different titles collide", func(t *testing.T) {
		a := Normalize("gh pr create --title 'Fix login bug' --body 'Fixes the auth issue'")
		b := Normalize("gh pr create --title 'Add feature X' --body 'New capability'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different PR numbers collide", func(t *testing.T) {
		a := Normalize("gh pr view 123")
		b := Normalize("gh pr view 456")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different repos collide", func(t *testing.T) {
		a := Normalize("gh pr list -R owner/repo1")
		b := Normalize("gh pr list -R owner/repo2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("gh pr view 123")
		subshell := Normalize("gh pr view $(get-pr-number)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
