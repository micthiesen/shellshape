---
name: review-approvals
description: Analyze recent approvals events for shellshape normalization gaps and prioritize or fix the resulting handler issues.
argument-hint: [hours (default: 24)]
---

# Review approvals

Query the recent event report with:

```bash
bash ${CLAUDE_SKILL_DIR}/scripts/query-events.sh $ARGUMENTS
```

The default window is 24 hours. Inspect singleton shapes, near-duplicates for one
executable, and high-cardinality executables. Literal paths, numbers, patterns,
headers, URLs, or other data often indicate under-normalization; distinct command
grammar that collapsed together indicates over-normalization.

For each credible issue, inspect the relevant handler and classify it as a missing
handler, handler bug, or shared token/preprocessing bug. Report the evidence,
category, and priority. High priority means common input or meaningful data leakage;
low priority means a rare form without practical impact.

If the user requested analysis, stop after the report. If they requested fixes,
use `add-handler` or `fix-handler` for clear issues in priority order. Ask only when
the intended shape or requested scope is materially ambiguous.
