#!/usr/bin/env bash
# go_lint_check.sh — Automated Go code quality checks
# Usage: ./go_lint_check.sh [project_dir]

set -euo pipefail

PROJECT_DIR="${1:-.}"
cd "$PROJECT_DIR"

ERRORS=0
WARNINGS=0

echo "══════════════════════════════════════════════"
echo "  🔍 Go Lint Check — $(basename "$(pwd)")"
echo "══════════════════════════════════════════════"
echo ""

# 1. go vet
echo "━━━ [1/4] go vet ━━━"
if go vet ./... 2>&1; then
    echo "✅ go vet: PASSED"
else
    echo "❌ go vet: FAILED"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# 2. staticcheck
echo "━━━ [2/4] staticcheck ━━━"
if command -v staticcheck &>/dev/null; then
    if staticcheck ./... 2>&1; then
        echo "✅ staticcheck: PASSED"
    else
        echo "❌ staticcheck: FAILED"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  staticcheck: NOT INSTALLED (go install honnef.co/go/tools/cmd/staticcheck@latest)"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 3. golangci-lint
echo "━━━ [3/4] golangci-lint ━━━"
if command -v golangci-lint &>/dev/null; then
    if golangci-lint run ./... 2>&1; then
        echo "✅ golangci-lint: PASSED"
    else
        echo "❌ golangci-lint: FAILED"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  golangci-lint: NOT INSTALLED (https://golangci-lint.run/welcome/install/)"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 4. go test -race
echo "━━━ [4/4] go test -race ━━━"
if go test -race -count=1 ./... 2>&1; then
    echo "✅ go test -race: PASSED"
else
    echo "❌ go test -race: FAILED"
    ERRORS=$((ERRORS + 1))
fi
echo ""

# Summary
echo "══════════════════════════════════════════════"
echo "  📊 Summary"
echo "══════════════════════════════════════════════"
echo "  Errors:   $ERRORS"
echo "  Warnings: $WARNINGS"

if [ "$ERRORS" -gt 0 ]; then
    echo "  Status:   ❌ FAILED"
    exit 1
else
    echo "  Status:   ✅ PASSED"
    exit 0
fi
