PROJECT_ROOT=$(shell pwd)

NAME=hub
TARGET_MAIN=
TARGET_EXT1=
TARGET_EXT2=

GOCMD=go
GOVER=1.20
CGO_ENABLED=0

GOBUILD=$(GOCMD) build -race
GOBUILD_ASAN=$(GOCMD) build -asan
GOBUILD_RELEASE=$(GOCMD) build -ldflags="-s -w -X main.version=${VERSION} -X main.revision=${REVISION} -X main.build=${BUILD}" -trimpath

GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=staticcheck
GOFMT=gofmt
GOALIGNMENT=fieldalignment
GOALIGNMENT_FIX=fieldalignment -fix

VERSION=$(shell git describe --tags --abbrev=0)
REVISION=$(shell git rev-parse --short HEAD)
BUILD=$(shell uname -s)_$(shell uname -r)_$(shell uname -m)
ARCH=$(shell uname -m)

.DEFAULT_GOAL := test
.PHONY: clean

init: _clean _envinit _gitconfiginit _modinit

_clean:
	-$(GOCLEAN) -modcache

_envinit:
	-go env -w GOPRIVATE=github.com/overlayos
	-go env -w GONOSUMDB=github.com/overlayos
	-go install honnef.co/go/tools/cmd/staticcheck@2023.1.7
	-go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest

_gitconfiginit:
	-git config pull.rebase true
	-git config --global core.editor vim

_modinit:
	$(GOCLEAN) -modcache
	- rm go.*
	go mod init github.com/overlayos/$(NAME)
	go mod tidy -go=$(GOVER)

_format:
	-$(GOFMT) -w ./

_lint:
	-$(GOLINT) ./...

_alignment:
	-$(GOALIGNMENT_FIX) .

test: _format _lint _alignment
	$(GOTEST) -v -race -shuffle on -cover -timeout 1m ./...

bench:
	$(GOTEST) -bench . -benchmem

profile:
	$(GOTEST) -cpuprofile badproto.pprof -memprofile mem.badproto.pprof -bench . -benchmem
