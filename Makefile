SHELL    := /bin/bash
MODULE   := $(shell sed -nr 's/^module ([a-z\-]+)$$/\1/p' go.mod)
GO_FILE  := src/$(MODULE).go
OUT_FILE := $(MODULE).bin

default:
	@go version

start:
	go run $(GO_FILE)

clean:
	go clean
	rm -f $(OUT_FILE)

# Build development binary.
build:
	go build -v -o $(OUT_FILE) $(GO_FILE)
	@ls -lh $(OUT_FILE)

# Build production binary.
release: clean
	@{\
		LD_FLAGS="-s -w" ;\
		CGO_ENABLED=0 go build -v -a \
			-ldflags "$$LD_FLAGS" \
			-o $(OUT_FILE) $(GO_FILE) ;\
	}
	@ls -lh $(OUT_FILE)
