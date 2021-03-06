VERSION?=$(shell git describe --tags --always --dirty)

all: build

build:
	go build -ldflags "-X main.version=$(VERSION)" ./cmd/nif

install: build
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/nif

dist: dist/nif_darwin_amd64 dist/nif_linux_amd64

dist/nif_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o dist/nif_darwin_amd64 -ldflags "-X main.version=$(VERSION)" ./cmd/nif

dist/nif_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o dist/nif_linux_amd64 -ldflags "-X main.version=$(VERSION)" ./cmd/nif

rel: dist
	hub release create -a dist $(VERSION)
