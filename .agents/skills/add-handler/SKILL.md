---
name: add-handler
description: Add or alias a shellshape handler for a command whose argument grammar needs command-specific normalization.
argument-hint: <executable-name>
---

# Add a shellshape handler

Add or alias a handler for the command or executable supplied by the user with the smallest grammar that preserves
command structure and collapses data.

## Decide the shape

First check whether an existing handler has the same flag and positional grammar.
If so, register the command as an alias and extend that handler's tests.

Otherwise, research the command's official documentation and representative
invocations. Identify:

- boolean flags and flags that consume paths, values, patterns, or numbers;
- positional roles and subcommands;
- tokens that are data and should collapse versus structure that must remain.

State the intended shapes as a brief progress update, then continue unless the
expected normalization is materially ambiguous.

## Implement from a failing test

Create `handlers/<command>_test.go` first. Cover common forms, flag arguments,
positionals, aliases or subcommands, and edge cases. Include:

- a collision test proving different data values produce the same shape;
- a safety test proving a `$(...)` token does not collapse like literal data.

Run the focused test and confirm that the new case fails for the expected reason.
Then create `handlers/<command>.go` or update the aliased handler.

Handlers register from `init()`. Pass
`shellshape.HandlerOptions{HasSubcommands: true}` when subcommands change the
grammar. In the token walk:

- call `shellshape.SplitRedirects(tokens)` first and append redirects last;
- check `shellshape.IsSubshellToken` before collapsing a positional;
- use `shellshape.IsFlagToken` and `shellshape.ClassifyToken` for shared rules;
- walk by index when flags consume the next token;
- use `ConsumeFlagArg`, `FlagCategory`, `MatchFlagCategory`, and
  `ConsumeFusedFlag` instead of duplicating flag parsing.

Read `handler.go` for the shared helpers and use a nearby handler with comparable
grammar as the local pattern.

## Verify

Format the changed Go files, then run:

```bash
go vet ./...
go test ./...
```

Fix regressions in the handler or a mistaken expected shape. Preserve the
subshell distinction in every path that collapses user input.
