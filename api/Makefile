VERSION=$(shell git describe --tags --candidates=1 --dirty)
BUILD_FLAGS=-ldflags="-X main.Version=$(VERSION)" -trimpath
SRC=$(shell find . -name '*.go') go.mod
.PHONY: nyla-api

nyla-api: $(SRC)
	go build -ldflags="-X main.Version=$(VERSION)" -o $@ .
