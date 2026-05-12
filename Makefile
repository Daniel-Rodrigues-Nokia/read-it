.DEFAULT_GOAL := build

VERSION ?= dev
TEST_TYPE ?= cypress

.PHONY: clean fmt vet build build-ci

clean:
				go clean
fmt: clean
				go fmt ./...
vet: fmt
				go vet ./...
build: vet
				go build -o bin/read-it-$(TEST_TYPE)
build-ci: vet
				go build -ldflags "-X main.version=$(VERSION) -X main.testType=$(TEST_TYPE)" -o bin/read-it-$(TEST_TYPE)
