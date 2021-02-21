SHELL   := /bin/bash
MODULE  := $(shell sed -nr 's/^module ([a-z\-]+)$$/\1/p' go.mod)
GO_FILE := src/program.go
DIST    := dist

default:
	go version

start:
	go run $(GO_FILE)

clean:
	go clean
	rm -rf $(DIST)/*

# Build development binary.
build:
	go build -v -o $(DIST)/$(MODULE) $(GO_FILE)
	ls -lh $(DIST)

# Build production binary.
release: clean
	@{\
		LD_FLAGS="-s -w" ;\
		CGO_ENABLED=0 go build -v -a \
			-ldflags "$$LD_FLAGS" \
			-o $(DIST)/$(MODULE) $(GO_FILE) ;\
	}
	ls -lh $(DIST)
