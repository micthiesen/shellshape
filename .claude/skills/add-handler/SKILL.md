---
name: add-handler
description: Add a new shellshape handler for a specific command (e.g. curl, awk, docker, jq)
argument-hint: <executable-name>
---

# Add a new shellshape handler

You are adding a handler for `$ARGUMENTS` to the shellshape normalizer. Follow this process exactly.

## Step 0: Check for an existing handler to alias

If the command is functionally identical to an existing handler (same flag grammar and positional semantics), just add it as an alias in that handler's `init()` instead of creating a new handler file. Add test cases for the new name to the existing test file and you're done.

## Step 1: Research the command

Understand how `$ARGUMENTS` is actually used. You need to know:

- Which flags take arguments (and what kind: paths, patterns, numbers, expressions, URLs)
- Which flags are boolean (no argument)
- What the positional arguments mean (is the first one special? are the rest files?)
- Whether it has subcommands (like `docker run` vs `docker build`)
- Common real-world invocations

Do this research by:

1. Run `curl -s 'cheat.sh/$ARGUMENTS?T'` to get a quick reference with common flags and examples
2. If you need more detail (e.g. obscure flags, exact argument semantics), look up the official man pages online or use other sources
3. Think about what tokens are **data** (should collapse) vs **structure** (should stay verbatim)

## Step 2: Design the normalization rules

Before writing any code, write out your plan:

- What placeholder does the "main" positional get? (e.g. `<pattern>` for grep, `<sed-expr>` for sed, `<filter>` for jq)
- Which flags consume the next token, and what placeholder does it get?
- Are there any flags that take numeric arguments (should become `N`)?
- Are there subcommands that change the argument grammar?
- What should the shape look like for 5-10 common invocations?

Share this plan with the user before proceeding.

## Step 3: Write the test file first (TDD)

Create `handler_$ARGUMENTS_test.go`. Follow the exact pattern from existing handler tests.

Structure your tests in these categories:

```go
func TestCommandName(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        // Basic usage
        // Flags with arguments
        // Multiple positionals
        // Edge cases
    }
    // ...

    // COLLISION TESTS: different data values → same shape
    t.Run("different X collide", func(t *testing.T) { ... })

    // SAFETY TEST: subshell must not collapse (one per handler, mandatory)
    t.Run("subshell not collapsed", func(t *testing.T) {
        benign := Normalize("$ARGUMENTS literal-arg")
        subshell := Normalize("$ARGUMENTS $(dangerous-command)")
        if benign == subshell {
            t.Error("subshell must produce different shape than literal")
        }
    })
}
```

Run the tests. They should fail. That's correct.

```bash
go test ./... -run TestCommandName -v
```

## Step 4: Implement the handler

Create `handler_$ARGUMENTS.go`. The handler self-registers via `init()`, so no other files need to be modified. Use this template:

```go
package shellshape

func init() {
    Register("$ARGUMENTS", handleCommandName)
}

func handleCommandName(subcommand string, tokens []string) []string {
    args, redirects := splitRedirects(tokens)

    var result []string
    // ... your logic here ...

    result = append(result, redirects...)
    return result
}
```

If the command has aliases (like grep/egrep/fgrep), register all of them in init():

```go
func init() {
    for _, name := range []string{"$ARGUMENTS", "alias1", "alias2"} {
        Register(name, handleCommandName)
    }
}
```

If the command has subcommands (like `docker run`, `git commit`), pass `HandlerOptions`:

```go
func init() {
    Register("$ARGUMENTS", handleCommandName, HandlerOptions{HasSubcommands: true})
}
```

Key rules:
- Always call `splitRedirects(tokens)` first, append redirects at end
- Always check `isSubshellToken(tok)` BEFORE collapsing any positional to a placeholder
- Use `isFlagToken(tok)` to distinguish flags from positionals
- Use `classifyToken(tok)` for positionals that should get generic classification (usually file paths)
- Walk tokens with an index variable (`i`), not `range`, when flags consume next args
- Use the shared utilities from `handler.go` for flag processing (read the file to see what's available):
  - `consumeFlagArg(tok, args, i, result, "<placeholder>")` to consume a flag and its next token (handles subshell preservation automatically)
  - `flagCategory` + `matchFlagCategory` to replace repetitive if/else chains when you have 3+ flag categories
  - `consumeFusedFlag(tok, categories)` to handle `--flag=value` syntax
- When a handler has 3+ flag categories (e.g. pathFlags, valFlags, numericFlags), define a `categories` slice and use `matchFlagCategory` in the loop instead of sequential if blocks

Reference `handler_grep.go` for a handler with flag-consuming arguments, or `handler_echo.go` for a simple positional-collapsing handler.

## Step 5: Run tests and iterate

```bash
go vet ./... && go test ./... -v
```

All tests must pass, including the existing ones (no regressions). If a test fails, fix the handler, not the test (unless you got the expected shape wrong in Step 3).

## Checklist before done

- [ ] Handler test file with 8+ test cases covering common usage
- [ ] At least one collision test (different data → same shape)
- [ ] Subshell safety test (mandatory, never skip)
- [ ] Handler implementation with `init()` self-registration
- [ ] `go vet ./... && go test ./...` passes clean
