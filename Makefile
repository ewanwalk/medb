VERSION=0.1.0
BUILD=$(shell git rev-parse HEAD)
TARGET=main

LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

.PHONY: all build fmt vet test

default: all

all: fmt vet test build

build:
	go build $(LDFLAGS) -o $(TARGET) cmd/medb/main.go

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)

test:
	go test -race -v $$(go list ./... | grep -v /vendor/)