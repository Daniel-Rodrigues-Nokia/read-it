.DEFAULT_GOAL := build

VERSION ?= dev
.PHONY: clean fmt vet build build-ci

clean:
				go clean
fmt: clean
				go fmt ./...
vet: fmt
				go vet ./...
build: vet
				go build -o bin/read-it 
build-ci: vet
				go build -ldflags "-X main.version=$(VERSION)" -o bin/read-it
