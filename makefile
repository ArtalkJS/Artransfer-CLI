PACKAGE_NAME := github.com/ArtalkJS/Artransfer-CLI
VERSION      ?= $(shell git describe --tags --abbrev=0)
COMMIT_HASH  := $(shell git rev-parse --short HEAD)
GO_VERSION   ?= 1.18.1

.PHONY: build
build:
	go build \
    	-ldflags "-s -w -X github.com/ArtalkJS/Artransfer-CLI/lib.Version=${VERSION} \
        -X github.com/ArtalkJS/Artransfer-CLI/lib.CommitHash=${COMMIT_HASH}" \
        -o bin/artransfer \
    	github.com/ArtalkJS/Artransfer-CLI

# https://github.com/goreleaser/goreleaser-cross-example
# https://github.com/goreleaser/goreleaser-cross
.PHONY: release
release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:v${GO_VERSION} \
		release --rm-dist --skip-validate
