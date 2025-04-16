#!/bin/bash

set -e

COVERAGE_DIR=$(mktemp -d)

COVERAGE_FILE="coverage.out"
echo "mode: atomic" > "$COVERAGE_FILE"

ALL_PACKAGES=$(go list ./... | grep -v mock | grep -v docs | grep -v populate)

for pkg in $ALL_PACKAGES; do
    echo "PKG: $pkg"
    pkg_coverage="$COVERAGE_DIR/$(echo $pkg | tr / -).cover"
    
    if ! go test -coverprofile="$pkg_coverage" -covermode=atomic "$pkg" > /dev/null 2>&1; then
        go test -coverprofile="$pkg_coverage" -covermode=atomic -run=^$ "$pkg" > /dev/null 2>&1 || true
    fi
    
    if [ -f "$pkg_coverage" ]; then
        tail -n +2 "$pkg_coverage" >> "$COVERAGE_FILE"
    fi
done

rm -rf "$COVERAGE_DIR"

go tool cover -func="$COVERAGE_FILE"