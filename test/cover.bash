#!/usr/bin/env bash
set -eux

# Clean previous coverage results if needed
rm -f *.covout coverage-report.txt

# Run the unit tests to retrieve coverage files
GOCACHE=off DO_COVERAGE=1 go test \
    -covermode=count \
    -coverprofile=unittest.covout \
    -coverpkg=github.com/netorcai/netorcai,github.com/netorcai/netorcai/cmd/netorcai \
    -v ..

# Run the integration tests to retrieve coverage files
GOCACHE=off DO_COVERAGE=1 go test -v . || true

# Merge all coverage files into one
gocovmerge *.covout > merged.covout

# Get a readable coverage report
gocov convert merged.covout | gocov report > coverage-report.txt
go tool cover -html=merged.covout -o coverage-report.html
