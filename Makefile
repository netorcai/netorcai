VERSION=`git describe --dirty`
LDFLAGS=-ldflags "-X main.version=${VERSION}"

netorcai:
	go build ${LDFLAGS} -o ./netorcai ./cmd/netorcai

netorcai.cover:
	go test -c -o ./netorcai.cover -covermode=count -coverpkg=./,./cmd/netorcai ./cmd/netorcai

rebuild-nocache:
	GOCACHE=off go build ${LDFLAGS} -o ./netorcai ./cmd/netorcai
	GOCACHE=off go test -c -o ./netorcai.cover -covermode=count -coverpkg=./,./cmd/netorcai ./cmd/netorcai

all: netorcai netorcai.cover

.PHONY: netorcai netorcai.cover
