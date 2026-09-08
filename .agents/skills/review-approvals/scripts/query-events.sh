#!/usr/bin/env bash
# Query the approvals event log for normalization review.
#
# Usage:
#   query-events.sh [hours]     # default: 24 hours
#
# Output sections:
#   1. Singleton shapes (hit only once) — likely normalization gaps
#   2. Near-duplicate shapes — similar shapes that might be collapsible
#   3. Raw commands for each singleton shape — for diagnosis
#
# The DB path matches the approvals system's default location.

set -euo pipefail

DB="${APPROVALS_DB:-$HOME/.local/share/claude-approvals/approvals.db}"
HOURS="${1:-24}"

if [[ ! -f "$DB" ]]; then
  echo "ERROR: database not found at $DB"
  echo "Set APPROVALS_DB to the correct path."
  exit 1
fi

CUTOFF=$(( $(date +%s) - HOURS * 3600 ))

echo "=== Approvals review: last ${HOURS}h (since $(date -r "$CUTOFF" '+%Y-%m-%d %H:%M')) ==="
echo ""

TOTAL=$(sqlite3 "$DB" "SELECT COUNT(*) FROM events WHERE timestamp >= $CUTOFF")
UNIQUE=$(sqlite3 "$DB" "SELECT COUNT(DISTINCT shape) FROM events WHERE timestamp >= $CUTOFF")
echo "Total events: $TOTAL"
echo "Unique shapes: $UNIQUE"
echo ""

echo "=== SINGLETON SHAPES (hit only once — possible normalization gaps) ==="
echo ""
sqlite3 -header -column "$DB" <<SQL
SELECT shape, raw_command, source, reason
FROM events
WHERE timestamp >= $CUTOFF
  AND shape IN (
    SELECT shape FROM events
    WHERE timestamp >= $CUTOFF
    GROUP BY shape
    HAVING COUNT(*) = 1
  )
ORDER BY shape;
SQL

echo ""
echo "=== NEAR-DUPLICATE SHAPES (shapes sharing an executable, sorted) ==="
echo ""
sqlite3 -header -column "$DB" <<SQL
WITH exe_shapes AS (
  SELECT DISTINCT shape,
    CASE
      WHEN shape LIKE '% %' THEN SUBSTR(shape, 1, INSTR(shape, ' ') - 1)
      ELSE shape
    END AS exe
  FROM events
  WHERE timestamp >= $CUTOFF
)
SELECT exe, GROUP_CONCAT(shape, char(10)) AS shapes
FROM exe_shapes
GROUP BY exe
HAVING COUNT(*) > 1
ORDER BY exe;
SQL

echo ""
echo "=== HIGH-CARDINALITY EXECUTABLES (many distinct shapes — normalization leaks?) ==="
echo ""
sqlite3 -header -column "$DB" <<SQL
WITH exe_shapes AS (
  SELECT DISTINCT shape,
    CASE
      WHEN shape LIKE '% %' THEN SUBSTR(shape, 1, INSTR(shape, ' ') - 1)
      ELSE shape
    END AS exe
  FROM events
  WHERE timestamp >= $CUTOFF
)
SELECT exe, COUNT(*) AS distinct_shapes
FROM exe_shapes
GROUP BY exe
HAVING COUNT(*) >= 3
ORDER BY distinct_shapes DESC;
SQL

echo ""
echo "=== ALL SHAPES WITH HIT COUNTS (descending) ==="
echo ""
sqlite3 -header -column "$DB" <<SQL
SELECT shape, COUNT(*) AS hits, source, reason
FROM events
WHERE timestamp >= $CUTOFF
GROUP BY shape
ORDER BY hits DESC
LIMIT 50;
SQL
