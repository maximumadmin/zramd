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
		CGO_ENABLED=0 go build -v -a -trimpath \
			-ldflags "$$LD_FLAGS" \
			-o $(OUT_FILE) $(GO_FILE) ;\
	}
	@ls -lh $(OUT_FILE)

# Run unit tests on all packages.
test:
	go test -v ./src/...

install:
	install -Dm755 $(OUT_FILE) "$(PREFIX)/usr/bin/$(MODULE)"
	install -Dm644 LICENSE -t "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	install -Dm644 extra/$(MODULE).default "$(PREFIX)/etc/default/$(MODULE)"
	install -Dm644 extra/$(MODULE).service -t "$(PREFIX)/usr/lib/systemd/system/"

uninstall:
	systemctl disable --now $(MODULE).service
	rm -f "$(PREFIX)/usr/lib/systemd/system/$(MODULE).service"
	rm -f "$(PREFIX)/etc/default/$(MODULE)"
	rm -rf "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	rm -f "$(PREFIX)/usr/bin/$(MODULE)"
