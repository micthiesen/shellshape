---
name: review-approvals
description: Review recent approvals data to find normalization gaps and improvements
argument-hint: [hours (default: 24)]
---

# Review approvals for normalization issues

Analyze the approvals event log to find commands that are normalizing poorly and could benefit from handler improvements.

## Step 1: Query the approvals database

Run the query script to pull recent events:

```bash
bash ${CLAUDE_SKILL_DIR}/scripts/query-events.sh $ARGUMENTS
```

If `$ARGUMENTS` is empty, it defaults to the last 24 hours.

## Step 2: Analyze the output

Look at each section of the report and identify issues:

### Singleton shapes
These are shapes that were seen only once. They *might* be normalization gaps where data leaked into the shape. For each singleton, ask:
- Does the shape contain a literal value that should have been collapsed? (a specific path that should be `<path>`, a number that should be `N`, a pattern that should be `<pattern>`)
- Is this genuinely a unique command, or did normalization fail to group it with similar commands?

### Near-duplicate shapes
These are multiple shapes for the same executable. Look for shapes that differ only in one token that should have been normalized away. For example:
- `git log --oneline -10` and `git log --oneline -20` should both be `git log --oneline N` (they are, but look for cases that slip through)
- `curl -H specific-header https://foo.com` producing different shapes per header value

### High-cardinality executables
If one executable has many distinct shapes, that's suspicious. Some commands genuinely have many forms (git), but others might be leaking data into shapes.

## Step 3: For each issue found, diagnose and classify

For each normalization problem, determine which category it falls into:

1. **Missing handler**: The command has no handler and needs one (e.g., `curl`, `jq`, `awk`). These commands have argument semantics that the generic classifier can't see.
2. **Handler bug**: The command has a handler but it's not covering a specific flag or positional pattern.
3. **Token classification bug**: The generic `classifyToken` function is failing to recognize a data pattern (new URL scheme, unusual path format, etc.)

## Step 4: Report findings

Present a summary table to the user:

| Shape / Command | Problem | Category | Priority |
|---|---|---|---|
| ... | ... | missing handler / handler bug / token bug | high/medium/low |

Priority guide:
- **High**: Command is common (high hit count or multiple singletons for same exe) and the leak is significant
- **Medium**: Uncommon command but clear normalization gap
- **Low**: Edge case, single occurrence, probably not worth a handler

## Step 5: Fix issues (with user approval)

For each issue the user wants to fix:
- **Missing handler**: Use `/add-handler <executable>` to create one
- **Handler or token bug**: Use `/fix-handler <the raw command that normalizes wrong>` to fix it

Ask the user which issues they want to address before invoking these skills.
