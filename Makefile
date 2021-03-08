SHELL    := /bin/bash
MODULE   := $(shell sed -nr 's/^module ([a-z\-]+)$$/\1/p' go.mod)
GO_FILE  := src/$(MODULE).go
ifeq ($(output),)
OUT_FILE := dist/$(MODULE).bin
else
OUT_FILE := $(output)
endif

default:
	@go version

start:
	go run $(GO_FILE)

clean:
	go clean
	rm -f "$(OUT_FILE)"

# Build development binary.
build:
	@{\
		if [[ "$(OUT_FILE)" == dist/* ]]; then \
			mkdir -p "$(OUT_FILE)" ;\
		fi ;\
	}
	go build -v -o $(OUT_FILE) $(GO_FILE)
	@ls -lh "$(OUT_FILE)"

# Build statically linked production binary.
release: clean
	@{\
		export GOFLAGS="-a -trimpath -ldflags=-w -ldflags=-s" ;\
		if [ "$${GOARCH}" != "arm" ]; then \
			export GOFLAGS="$${GOFLAGS} -buildmode=pie" ;\
		fi ;\
		CGO_ENABLED=0 go build -o "$(OUT_FILE)" $(GO_FILE) ;\
	}
	@ls -lh "$(OUT_FILE)"

# Build dinamically linked production binary.
release-dynamic: clean
	@{\
		export CGO_CPPFLAGS="$${CPPFLAGS}" ;\
		export CGO_CFLAGS="$${CFLAGS}" ;\
		export CGO_CXXFLAGS="$${CXXFLAGS}" ;\
		export CGO_LDFLAGS="$${LDFLAGS}" ;\
		export GOFLAGS="-a -trimpath -ldflags=-linkmode=external -ldflags=-w -ldflags=-s" ;\
		if [ "$${GOARCH}" != "arm" ]; then \
			export GOFLAGS="$${GOFLAGS} -buildmode=pie" ;\
		fi ;\
		go build -o "$(OUT_FILE)" $(GO_FILE) ;\
	}
	@ls -lh "$(OUT_FILE)"

# Run unit tests on all packages.
test:
	go test -v ./src/...

install:
	install -Dm755 "$(OUT_FILE)" "$(PREFIX)/usr/bin/$(MODULE)"
	install -Dm644 LICENSE -t "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	install -Dm644 extra/$(MODULE).default "$(PREFIX)/etc/default/$(MODULE)"
	install -Dm644 extra/$(MODULE).service -t "$(PREFIX)/usr/lib/systemd/system/"

uninstall:
	systemctl disable --now $(MODULE).service
	rm -f "$(PREFIX)/usr/lib/systemd/system/$(MODULE).service"
	rm -f "$(PREFIX)/etc/default/$(MODULE)"
	rm -rf "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	rm -f "$(PREFIX)/usr/bin/$(MODULE)"
