SHELL    := /bin/bash
MODULE   := $(shell sed -nr 's/^module ([a-z\-]+)$$/\1/p' go.mod)
GO_FILE  := src/$(MODULE).go
ifeq ($(output),)
OUT_FILE := dist/$(MODULE).bin
else
OUT_FILE := $(output)
endif

default:
	@{\
		set -e ;\
		os_release_id=$$(grep -E '^ID=' /etc/os-release | sed 's/ID=//' || true) ;\
		if [ "$$os_release_id" = "arch" ]; then \
			make --no-print-directory release-dynamic ;\
		else \
			make --no-print-directory release-static ;\
		fi ;\
	}

start:
	go run $(GO_FILE)

clean:
	go clean
	rm -rf dist/*
	rm -f "$(OUT_FILE)"

# Build development binary
build:
	go build -v -o $(OUT_FILE) $(GO_FILE)
	@ls -lh "$(OUT_FILE)"

# Build statically linked production binary
release-static:
	@echo "Building static binary (GOARCH: $(GOARCH) GOARM: $(GOARM))..."
	@{\
		set -e ;\
		if [ -z "$${skip_clean}" ]; then make --no-print-directory clean; fi ;\
		export GOFLAGS="-a -trimpath -ldflags=-w -ldflags=-s" ;\
		if [ "$${GOARCH}" != "arm" ]; then \
			export GOFLAGS="$${GOFLAGS} -buildmode=pie" ;\
		fi ;\
		CGO_ENABLED=0 go build -o "$(OUT_FILE)" $(GO_FILE) ;\
	}
	@make --no-print-directory postbuild

# Build dinamically linked production binary
release-dynamic:
	@echo "Building dynamic binary (GOARCH: $(GOARCH) GOARM: $(GOARM))..."
	@{\
		set -e ;\
		if [ -z "$${skip_clean}" ]; then make --no-print-directory clean; fi ;\
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
	@make --no-print-directory postbuild

postbuild:
	@{\
		set -e ;\
		if [ ! -z "$${make_tgz}" ]; then \
			tgz_file="$(OUT_FILE).tar.gz" ;\
			echo "Creating \"$${tgz_file}\"..." ;\
			tar -C "$$(dirname "$(OUT_FILE)")" \
				-cz -f "$$tgz_file" \
				"$$(basename "$(OUT_FILE)")" ;\
		fi ;\
		if [ ! -z "$${make_deb}" ]; then \
			echo "Creating deb ($${DEB_ARCH}) file..." ;\
			CONFIG_FILE=extra/debian.yml \
				ARCH=$${DEB_ARCH} \
				PREFIX="$${PREFIX}" \
				VERSION=$${VERSION} \
				RELEASE=$${RELEASE} \
				./scripts/mkdeb.py ;\
		fi ;\
	}
	@ls -lh "$(OUT_FILE)"*

# Run unit tests on all packages
test:
	go test -v ./src/...

install:
	@{\
		set -e ;\
		bin_file="$${BIN_FILE:-$(OUT_FILE)}" ;\
		install -Dm755 "$${bin_file}" "$(PREFIX)/usr/bin/$(MODULE)" ;\
	}
	install -Dm644 LICENSE -t "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	install -Dm644 extra/$(MODULE).default "$(PREFIX)/etc/default/$(MODULE)"
	install -Dm644 extra/$(MODULE).service -t "$(PREFIX)/usr/lib/systemd/system/"

uninstall:
	@{\
		if [ -f "$(PREFIX)/usr/lib/systemd/system/$(MODULE).service" ]; then \
			systemctl disable --now $(MODULE).service ;\
		fi ;\
	}
	rm -f "$(PREFIX)/usr/lib/systemd/system/$(MODULE).service"
	rm -f "$(PREFIX)/etc/default/$(MODULE)"
	rm -rf "$(PREFIX)/usr/share/licenses/$(MODULE)/"
	rm -f "$(PREFIX)/usr/bin/$(MODULE)"
