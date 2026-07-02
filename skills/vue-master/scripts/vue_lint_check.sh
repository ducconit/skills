#!/usr/bin/env bash
# vue_lint_check.sh — Automated Vue.js code quality checks
# Usage: ./vue_lint_check.sh [project_dir]

set -euo pipefail

PROJECT_DIR="${1:-.}"
cd "$PROJECT_DIR"

ERRORS=0
WARNINGS=0

echo "══════════════════════════════════════════════"
echo "  🔍 Vue.js Lint Check — $(basename "$(pwd)")"
echo "══════════════════════════════════════════════"
echo ""

# Check if node_modules exists
if [ ! -d "node_modules" ]; then
    echo "⚠️  node_modules not found. Run 'bun install' first."
    exit 1
fi

# 1. ESLint
echo "━━━ [1/3] ESLint ━━━"
if [ -f ".eslintrc.js" ] || [ -f ".eslintrc.cjs" ] || [ -f ".eslintrc.json" ] || [ -f "eslint.config.js" ] || [ -f "eslint.config.mjs" ]; then
    if npx eslint --ext .vue,.ts,.tsx,.js src/ 2>&1; then
        echo "✅ ESLint: PASSED"
    else
        echo "❌ ESLint: FAILED"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  ESLint config not found. Consider adding ESLint."
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 2. TypeScript Check (vue-tsc)
echo "━━━ [2/3] vue-tsc TypeScript Check ━━━"
if command -v npx &>/dev/null && npx vue-tsc --version &>/dev/null 2>&1; then
    if npx vue-tsc --noEmit 2>&1; then
        echo "✅ vue-tsc: PASSED"
    else
        echo "❌ vue-tsc: FAILED"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  vue-tsc not found. Install: npm install -D vue-tsc"
    WARNINGS=$((WARNINGS + 1))
fi
echo ""

# 3. Prettier
echo "━━━ [3/3] Prettier Format Check ━━━"
if [ -f ".prettierrc" ] || [ -f ".prettierrc.json" ] || [ -f ".prettierrc.js" ] || [ -f "prettier.config.js" ]; then
    if npx prettier --check 'src/**/*.{vue,ts,tsx,js,css}' 2>&1; then
        echo "✅ Prettier: PASSED"
    else
        echo "❌ Prettier: FAILED (run 'npx prettier --write' to fix)"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo "⚠️  Prettier config not found. Consider adding Prettier."
    WARNINGS=$((WARNINGS + 1))
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
