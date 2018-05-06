VERSION=`git describe --dirty`
LDFLAGS=-ldflags "-X main.version=${VERSION}"

netorcai: setup
	go build ${LDFLAGS} -o ./netorcai ./cmd/netorcai

netorcai.cover: setup
	go test -c -o ./netorcai.cover -covermode=count -coverpkg=./,./cmd/netorcai ./cmd/netorcai

rebuild-nocache: setup
	GOCACHE=off go build ${LDFLAGS} -o ./netorcai ./cmd/netorcai
	GOCACHE=off go test -c -o ./netorcai.cover -covermode=count -coverpkg=./,./cmd/netorcai ./cmd/netorcai

setup:
	go get ./
	go get ./cmd/netorcai

all: netorcai netorcai.cover

.PHONY: netorcai netorcai.cover
