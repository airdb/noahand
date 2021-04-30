BINARY := $(shell basename "$(PWD)")
VERSION:=$(shell git describe --dirty --always)
BUILD := $(shell git rev-parse HEAD)

SYSTEM:=
#LDFLAGS=-ldflags
#LDFLAGS += "-X=github.com/airdb/adb/internal/adblib.Version=$(VERSION) \
#            -X=github.com/airdb/adb/internal/adblib.Build=$(BUILD) \
#            -X=github.com/airdb/adb/internal/adblib.BuildTime=$(shell date +%s)"

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
	CGO_ENABLED=0 GOOS=$(os) GOARCH=amd64 go build $(LDFLAGS) -o release/$(BINARY)-$(os)

.PHONY: release
release: windows linux darwin
