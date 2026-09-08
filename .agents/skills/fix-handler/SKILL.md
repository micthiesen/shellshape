---
name: fix-handler
description: Diagnose and fix a shellshape command that normalizes to the wrong structural shape.
argument-hint: <command-that-normalizes-wrong>
---

# Fix a shellshape handler

Reproduce `Normalize` with the command supplied by the user and determine the expected structural shape.
Trace the fault to the executable handler, `ClassifyToken`, or preprocessing.
Infer the expected result from established rules; ask only if distinct reasonable
shapes would materially change behavior. Report the diagnosis briefly, then continue.

Add a focused failing test in the owning test file before changing production
code. Make the smallest fix that addresses the demonstrated grammar and related
cases using the same path. Preserve these invariants:

- structurally different commands do not collide;
- different data values do collide when their role is the same;
- `$(...)` tokens never collapse to an ordinary data placeholder;
- redirects and flag arguments retain their established treatment.

Run the focused test while iterating, then format and verify the repository:

```bash
go vet ./...
go test ./...
```
