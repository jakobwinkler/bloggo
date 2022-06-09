VERSION:=$(shell git describe --tags --always --dirty)

all:
	go build -ldflags "-X github.com/jakobwinkler/bloggo/util.Version=${VERSION}"
