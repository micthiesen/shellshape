---
name: fix-handler
description: Fix or improve a shellshape handler when a command normalizes incorrectly
argument-hint: <command-that-normalizes-wrong>
---

# Fix a shellshape handler

The user has found a command that normalizes incorrectly. The input command is:

```
$ARGUMENTS
```

## Step 1: Reproduce and diagnose

First, see what we currently produce:

```go
Normalize("$ARGUMENTS")
```

Run this in a test or via the CLI to see the actual output. Then determine:

- What shape does it produce now?
- What shape *should* it produce? (Ask the user if unclear)
- Is the problem in a per-executable handler, the generic token classifier (`classifyToken`), or the preprocessing pipeline (shlex, heredocs, comments, splitting)?

To figure out where the bug is:
1. Check if the executable has a handler (look for `handler_<exe>.go` and its `init()` registration)
2. If it does, read the handler and trace through the token walk mentally
3. If it doesn't, the issue is in `classifyToken` in `token.go` or in preprocessing (`preprocess.go`, `split.go`, `shlex.go`)

Share your diagnosis with the user before proceeding.

## Step 2: Write a failing test first

Add the failing case to the appropriate test file:

- Handler bug → `handler_<exe>_test.go`
- Token classification bug → `token_test.go`
- Pipeline bug → `normalize_test.go`

Run it and confirm it fails:

```bash
go test ./... -run TestName/case_name -v
```

## Step 3: Fix it

Make the minimal change needed. Watch out for:

- **Over-normalization**: collapsing tokens that should stay distinct (two structurally different commands producing the same shape)
- **Under-normalization**: leaving data tokens verbatim that should collapse (two data-different commands producing different shapes)
- **Subshell safety**: never collapse a `$(...)` token to a data placeholder. If your fix touches positional handling, verify the subshell safety test still passes
- **Regressions**: a fix for one command can break another. Run the full suite, not just your new test

## Step 4: Check for related cases

If the bug was in a handler, think about whether similar commands have the same problem. For example, if `grep -A` was mishandling its numeric arg, check `grep -B` and `grep -C` too.

If the bug was in `classifyToken`, think about what other inputs could hit the same code path.

Add test cases for any related patterns you find.

## Step 5: Verify

```bash
go vet ./... && go test ./...
```

All tests must pass, including the new ones and all existing ones.
