#!/usr/bin/env bash
# go_complexity_check.sh — Cyclomatic complexity analysis
# Usage: ./go_complexity_check.sh [project_dir] [threshold]

set -euo pipefail

PROJECT_DIR="${1:-.}"
THRESHOLD="${2:-15}"
cd "$PROJECT_DIR"

echo "══════════════════════════════════════════════"
echo "  🧮 Go Complexity Check — $(basename "$(pwd)")"
echo "  Threshold: $THRESHOLD"
echo "══════════════════════════════════════════════"
echo ""

TOOL=""
FOUND_ISSUES=0

# Try gocyclo first, then gocognit
if command -v gocyclo &>/dev/null; then
    TOOL="gocyclo"
elif command -v gocognit &>/dev/null; then
    TOOL="gocognit"
fi

if [ -z "$TOOL" ]; then
    echo "⚠️  Neither gocyclo nor gocognit found."
    echo "   Install one:"
    echo "   go install github.com/fzipp/gocyclo/cmd/gocyclo@latest"
    echo "   go install github.com/uudashr/gocognit/cmd/gocognit@latest"
    echo ""
    echo "Falling back to basic analysis..."
    echo ""

    # Fallback: count lines per function using grep
    echo "━━━ Functions with many lines (rough estimate) ━━━"
    find . -name '*.go' -not -path './vendor/*' -not -name '*_test.go' | while read -r file; do
        awk '/^func /{name=$0; lines=0; next} /^}/{if(lines>50) printf "⚠️  %s:%s (%d lines)\n", FILENAME, name, lines; lines=0; next} {lines++}' "$file"
    done
    exit 0
fi

echo "Using: $TOOL"
echo ""

# Run complexity check
echo "━━━ Functions exceeding threshold ($THRESHOLD) ━━━"
RESULT=$($TOOL -over "$THRESHOLD" ./... 2>&1 || true)

if [ -z "$RESULT" ]; then
    echo "✅ All functions are within complexity threshold"
else
    echo "$RESULT"
    FOUND_ISSUES=$(echo "$RESULT" | wc -l)
    echo ""
    echo "━━━ Refactoring Suggestions ━━━"
    echo "For each complex function above, consider:"
    echo "  1. Extract helper functions for distinct logic blocks"
    echo "  2. Use early returns to reduce nesting"
    echo "  3. Replace switch/if chains with strategy pattern or maps"
    echo "  4. Move validation logic to separate validator functions"
    echo "  5. Use table-driven approach for repetitive conditionals"
fi

echo ""
echo "══════════════════════════════════════════════"
echo "  📊 Summary"
echo "══════════════════════════════════════════════"
echo "  Complex functions found: $FOUND_ISSUES"

if [ "$FOUND_ISSUES" -gt 0 ]; then
    echo "  Status: ⚠️  NEEDS ATTENTION"
    exit 1
else
    echo "  Status: ✅ CLEAN"
    exit 0
fi
