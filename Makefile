VERSION=`git describe --dirty`
LDFLAGS=-ldflags "-X main.version=${VERSION}"

netorcai:
	go build ${LDFLAGS} -o ./netorcai ./

netorcai.cover:
	go test -c -o ./netorcai.cover -covermode=count -coverpkg=./ ./

all: netorcai netorcai.cover

.PHONY: netorcai netorcai.cover
