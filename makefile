VERSION      ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH  := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/Artransfer-CLI/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/Artransfer-CLI/lib.CommitHash=${COMMIT_HASH}" \
        -o bin/artransfer \
    	github.com/ArtalkJS/Artransfer-CLI
