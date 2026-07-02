#!/usr/bin/env bash
# vue_bundle_check.sh — Bundle size analysis for Vue.js projects
# Usage: ./vue_bundle_check.sh [project_dir]

set -euo pipefail

PROJECT_DIR="${1:-.}"
cd "$PROJECT_DIR"

echo "══════════════════════════════════════════════"
echo "  📦 Vue.js Bundle Analysis — $(basename "$(pwd)")"
echo "══════════════════════════════════════════════"
echo ""

# Check if package.json exists
if [ ! -f "package.json" ]; then
    echo "❌ package.json not found"
    exit 1
fi

WARNINGS=0

# 1. Check for large dependencies
echo "━━━ [1/3] Large Dependencies Check ━━━"
echo "Checking dependencies > 100KB..."
echo ""

# Known large dependencies to watch for
LARGE_DEPS=("moment" "lodash" "luxon" "date-fns" "chart.js" "three" "d3" "xlsx" "pdf-lib")

for dep in "${LARGE_DEPS[@]}"; do
    if grep -q "\"$dep\"" package.json 2>/dev/null; then
        echo "⚠️  Found '$dep' — consider:"
        case "$dep" in
            "moment")
                echo "   → Use 'dayjs' (2KB) instead of moment (300KB+)"
                WARNINGS=$((WARNINGS + 1))
                ;;
            "lodash")
                echo "   → Use 'lodash-es' for tree-shaking, or individual imports"
                echo "   → import { debounce } from 'lodash-es'"
                WARNINGS=$((WARNINGS + 1))
                ;;
            *)
                echo "   → Verify this is needed, consider lazy loading"
                ;;
        esac
    fi
done

if [ "$WARNINGS" -eq 0 ]; then
    echo "✅ No known large dependency issues found"
fi
echo ""

# 2. Check for lazy loading opportunities
echo "━━━ [2/3] Lazy Loading Opportunities ━━━"

# Check routes for lazy loading
ROUTE_FILES=$(find src -name "router*" -o -name "routes*" 2>/dev/null | head -10)
if [ -n "$ROUTE_FILES" ]; then
    echo "Route files found:"
    for file in $ROUTE_FILES; do
        STATIC_IMPORTS=$(grep -c "^import.*from" "$file" 2>/dev/null || echo "0")
        LAZY_IMPORTS=$(grep -c "() => import\|defineAsyncComponent\|lazy(" "$file" 2>/dev/null || echo "0")
        echo "  📄 $file"
        echo "     Static imports: $STATIC_IMPORTS"
        echo "     Lazy imports:   $LAZY_IMPORTS"
        if [ "$STATIC_IMPORTS" -gt 3 ] && [ "$LAZY_IMPORTS" -eq 0 ]; then
            echo "     ⚠️  Consider lazy loading route components"
            WARNINGS=$((WARNINGS + 1))
        fi
    done
else
    echo "  No route files found for analysis"
fi
echo ""

# 3. Check for tree-shaking issues
echo "━━━ [3/3] Tree-Shaking Checks ━━━"

# Check for barrel imports that hurt tree-shaking
BARREL_IMPORTS=$(grep -r "import \* as" src/ --include="*.ts" --include="*.vue" 2>/dev/null | head -10)
if [ -n "$BARREL_IMPORTS" ]; then
    echo "⚠️  Found wildcard imports (harmful for tree-shaking):"
    echo "$BARREL_IMPORTS" | head -5
    WARNINGS=$((WARNINGS + 1))
else
    echo "✅ No wildcard imports found"
fi
echo ""

# Check for sideEffects in package.json
if grep -q '"sideEffects"' package.json 2>/dev/null; then
    echo "✅ 'sideEffects' field configured in package.json"
else
    echo "ℹ️  Consider adding 'sideEffects: false' to package.json for better tree-shaking"
fi
echo ""

# Build size check
echo "━━━ Build Size (if dist exists) ━━━"
if [ -d "dist" ]; then
    TOTAL_SIZE=$(du -sh dist/ | cut -f1)
    JS_SIZE=$(find dist -name "*.js" -exec du -ch {} + 2>/dev/null | tail -1 | cut -f1)
    CSS_SIZE=$(find dist -name "*.css" -exec du -ch {} + 2>/dev/null | tail -1 | cut -f1)

    echo "  Total dist size: $TOTAL_SIZE"
    echo "  JS size:         $JS_SIZE"
    echo "  CSS size:        $CSS_SIZE"

    # Check for files > 250KB
    echo ""
    echo "  Large files (> 250KB):"
    LARGE_FILES=$(find dist -name "*.js" -size +250k 2>/dev/null)
    if [ -n "$LARGE_FILES" ]; then
        echo "$LARGE_FILES" | while read -r f; do
            SIZE=$(du -h "$f" | cut -f1)
            echo "    ⚠️  $f ($SIZE)"
        done
        WARNINGS=$((WARNINGS + 1))
    else
        echo "    ✅ No JS files > 250KB"
    fi
else
    echo "  ℹ️  dist/ not found. Run 'npm run build' first for size analysis."
fi
echo ""

# Summary
echo "══════════════════════════════════════════════"
echo "  📊 Summary"
echo "══════════════════════════════════════════════"
echo "  Warnings: $WARNINGS"

if [ "$WARNINGS" -gt 0 ]; then
    echo "  Status:   ⚠️  NEEDS ATTENTION"
else
    echo "  Status:   ✅ CLEAN"
fi
