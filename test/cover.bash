#!/usr/bin/env bash
set -eux

# Create bats files with coverage
./generate_cover_bats_files.py

# Clean previous coverage results if needed
rm -f *.covout coverage-report.txt

# Run the tests, so coverage files can be obtained
GOCACHE=off DO_COVERAGE=1 go test .

# Merge all coverage files into one
gocovmerge *.covout > merged.covout

# Get a readable coverage report
gocov convert merged.covout | gocov report > coverage-report.txt
go tool cover -html=merged.covout -o coverage-report.html
