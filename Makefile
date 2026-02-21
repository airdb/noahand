SERIVCE := noahand
REPO := github.com/airdb/noahand

BINARY := $(shell basename "$(PWD)")
VERSION:=$(shell git describe --dirty --always)
BUILD := $(shell git rev-parse HEAD)
BUILDTIME := $(shell date +%s)"

LDFLAGS=-ldflags
LDFLAGS += "-X=$(REPO)/internal/version.Repo=$(REPO) \
            -X=$(REPO)/internal/version.Version=$(VERSION) \
            -X=$(REPO)/internal/version.Build=$(BUILD) \
            -X=$(REPO)/internal/version.BuildTime=${BUILDTIME}


all: help

help: ## Show help messages
	@echo "Container - ${SERVICE} "
	@echo
	@echo "Usage:\tmake COMMAND"
	@echo
	@echo "Commands:"
	@sed -n '/##/s/\(.*\):.*##/  \1#/p' ${MAKEFILE_LIST} | grep -v "MAKEFILE_LIST" | column -t -c 2 -s '#'

.PHONY: test

test:
	go test -v ./...

build:
	$(SYSTEM) go build $(LDFLAGS) -o main main.go
	#$(SYSTEM) GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY) cmd/cli/main.go
	#$(SYSTEM) GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -o plugins/plugin_greeter.so  plugins/greeter.go

PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build $(LDFLAGS) -o release/$(BINARY)-$(os) cmd/cli/main.go

.PHONY: release
release: linux darwin
	tar czvf release/$(BINARY)_latest.zip release/$(BINARY)-*
