PROJECT          := github.com/pulumi/go-yaml
TESTPARALLELISM := 10

WORKING_DIR     := $(shell pwd)

.PHONY: ensure build test

ensure::
	@echo "go mod tidy"; go mod tidy

build::
	@echo "go build ./..."; go build ./...

test::
	go test ./... -parallel ${TESTPARALLELISM} -timeout 10m
