#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m'

rm -f coverage.out coverage.html

touch coverage.tmp

PACKAGES=$(go list ./... | grep -v "mock" | grep -v "/gen/" | grep -v "/pb$" | grep -v "docs")

for pkg in $PACKAGES; do
    echo -e "${GREEN}Testing: ${NC}$pkg"
    
    go test -coverprofile=profile.out -covermode=atomic $pkg
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Tests failed for package: ${NC}$pkg"
        exit 1
    fi
    
    if [ -f profile.out ]; then
        grep -v "_easyjson.go" profile.out > profile_filtered.out 2>/dev/null || cp profile.out profile_filtered.out
        
        if [ ! -s coverage.tmp ]; then
            cat profile_filtered.out > coverage.tmp
        else
            tail -n +2 profile_filtered.out >> coverage.tmp
        fi
        rm profile.out profile_filtered.out
    fi
done

mv coverage.tmp coverage.out

go tool cover -html=coverage.out -o coverage.html

total_coverage=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}')
echo "$total_coverage" > coverage_percent.txt
echo -e "${GREEN}Total test coverage: ${NC}$total_coverage"

echo -e "${GREEN}Testing completed successfully!${NC}"
echo -e "Coverage report available at: ${YELLOW}coverage.html${NC}"

