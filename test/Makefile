default:
	@echo "No default target."

setup:
	@go get github.com/stretchr/testify/assert

clean:
	rm -f *.covout

cover: setup clean
	./cover.bash

test: setup clean
	@unset DO_COVERAGE || true
	GOCACHE=off go test -v .
