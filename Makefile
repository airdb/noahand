BINARY := $(shell basename "$(PWD)")
VERSION:=$(shell git describe --dirty --always)
BUILD := $(shell git rev-parse HEAD)

SYSTEM:=
REPO := github.com/airdb/noah

LDFLAGS=-ldflags
LDFLAGS += "-X=$(REPO)/internal/version.Repo=$(REPO) \
            -X=$(REPO)/internal/version.Version=$(VERSION) \
            -X=$(REPO)/internal/version.Build=$(BUILD) \
            -X=$(REPO)/internal/version.BuildTime=$(shell date +%s)"

.PHONY: test

all: build

test:
	go test -v ./...

build:
	$(SYSTEM) GOARCH=amd64 go build $(LDFLAGS) -o main main.go
	$(SYSTEM) GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY) cmd/cli/main.go
	#$(SYSTEM) GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=plugin -o plugins/plugin_greeter.so  plugins/greeter.go

PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build $(LDFLAGS) -o release/$(BINARY)-$(os) cmd/cli/main.go

.PHONY: release
release: linux darwin
	tar czvf release/$(BINARY)_latest.zip service release/$(BINARY)-*

installer:
	tar czvf release/$(BINARY)_latest.zip service release/$(BINARY)-*
	cat service/self_extracting_script.sh release/noah_latest.zip  > /tmp/install.sh

arm:
	GOOS=linux GOARCH=arm GOARM=7 go build $(LDFLAGS) -o $(BINARY) cmd/cli/main.go
