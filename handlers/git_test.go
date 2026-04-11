package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGit(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// commit
		{"commit with message", "git commit -m 'fix login bug'", "git commit -m <str>"},
		{"commit with long message flag", "git commit --message 'update readme'", "git commit --message <str>"},
		{"commit with author", "git commit --author 'John Doe <john@example.com>'", "git commit --author <val>"},
		{"commit amend no edit", "git commit --amend --no-edit", "git commit --amend --no-edit"},
		{"commit specific files", "git commit -m 'fix' src/main.go src/util.go", "git commit -m <str> <path>+"},
		{"commit all with message", "git commit -am 'quick fix'", "git commit -am <str>"},

		// add
		{"add single file", "git add src/main.go", "git add <path>"},
		{"add multiple files", "git add file1.js file2.js file3.js", "git add <path>+"},
		{"add all", "git add -A", "git add -A"},
		{"add patch", "git add -p src/main.go", "git add -p <path>"},

		// push / pull / fetch
		{"push simple", "git push", "git push"},
		{"push remote branch", "git push origin main", "git push origin main"},
		{"push force", "git push --force-with-lease origin feature/foo", "git push --force-with-lease origin <path>"},
		{"pull", "git pull origin main", "git pull origin main"},
		{"pull rebase", "git pull --rebase", "git pull --rebase"},
		{"fetch", "git fetch --all", "git fetch --all"},
		{"fetch prune", "git fetch --prune origin", "git fetch --prune origin"},

		// log
		{"log simple", "git log --oneline", "git log --oneline"},
		{"log with count", "git log -n 10 --oneline", "git log -n N --oneline"},
		{"log with format", "git log --format '%H %s'", "git log --format <val>"},
		{"log with pretty", "git log --pretty=oneline", "git log --pretty=<val>"},
		{"log with author filter", "git log --author john", "git log --author <val>"},
		{"log since", "git log --since '2024-01-01'", "git log --since <val>"},

		// diff
		{"diff simple", "git diff", "git diff"},
		{"diff staged", "git diff --staged", "git diff --staged"},
		{"diff with file", "git diff HEAD src/main.go", "git diff HEAD <path>"},
		{"diff rev range", "git diff main...feature", "git diff <range>"},
		{"diff context lines", "git diff -U 5", "git diff -U N"},

		// clone
		{"clone https", "git clone https://github.com/user/repo.git", "git clone <https-uri>"},
		{"clone ssh", "git clone git@github.com:user/repo.git", "git clone <git-uri>"},
		{"clone with dir", "git clone https://github.com/user/repo.git mydir", "git clone <https-uri> mydir"},
		{"clone depth", "git clone --depth 1 https://github.com/user/repo.git", "git clone --depth N <https-uri>"},

		// checkout / switch / branch
		{"checkout branch", "git checkout main", "git checkout main"},
		{"checkout new branch", "git checkout -b new-feature", "git checkout -b new-feature"},
		{"switch branch", "git switch main", "git switch main"},
		{"branch list", "git branch", "git branch"},
		{"branch delete", "git branch -d old-branch", "git branch -d old-branch"},

		// rebase / merge / reset
		{"rebase onto", "git rebase --onto main feature", "git rebase --onto <val> feature"},
		{"rebase interactive", "git rebase -i HEAD~3", "git rebase -i <rev>"},
		{"merge branch", "git merge feature-branch", "git merge feature-branch"},
		{"reset soft", "git reset --soft HEAD~1", "git reset --soft <rev>"},
		{"reset file", "git reset HEAD src/main.go", "git reset HEAD <path>"},

		// stash
		{"stash push", "git stash push -m 'work in progress'", "git stash push -m <str>"},
		{"stash pop", "git stash pop", "git stash pop"},
		{"stash drop index", "git stash drop stash@{2}", "git stash drop stash@{N}"},

		// tag
		{"tag annotated", "git tag -a v1.0.0 -m 'release 1.0'", "git tag -a v1.0.0 -m <str>"},
		{"tag delete", "git tag -d v1.0.0", "git tag -d v1.0.0"},

		// remote
		{"remote add", "git remote add origin git@github.com:user/repo.git", "git remote add origin <git-uri>"},
		{"remote remove", "git remote remove upstream", "git remote remove upstream"},

		// config
		{"config set", "git config user.name 'John Doe'", "git config <dotted-id> <val>"},
		{"config get", "git config --global user.email", "git config --global <dotted-id>"},

		// show
		{"show commit", "git show abc1234", "git show <hash>"},
		{"show head", "git show HEAD", "git show HEAD"},

		// cherry-pick / revert
		{"cherry-pick", "git cherry-pick abc1234def", "git cherry-pick <hash>"},
		{"revert", "git revert HEAD~2", "git revert <rev>"},

		// rm
		{"rm file", "git rm src/old.go", "git rm <path>"},
		{"rm cached", "git rm --cached secrets.env", "git rm --cached <path>"},

		// global flags
		{"global -C flag", "git -C /path/to/repo status", "git -C <path> status"},

		// redirects
		{"log with redirect", "git log --oneline > output.txt", "git log --oneline > <path>"},
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
	t.Run("different commit messages collide", func(t *testing.T) {
		a := shellshape.Normalize("git commit -m 'fix login bug'")
		b := shellshape.Normalize("git commit -m 'update readme section'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different clone URLs collide", func(t *testing.T) {
		a := shellshape.Normalize("git clone https://github.com/user/repo1.git")
		b := shellshape.Normalize("git clone https://github.com/other/repo2.git")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different file paths in add collide", func(t *testing.T) {
		a := shellshape.Normalize("git add src/main.go")
		b := shellshape.Normalize("git add lib/util.ts")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("git add literal-arg")
		subshell := shellshape.Normalize("git add $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
