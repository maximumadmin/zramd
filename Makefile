SHELL    := /bin/bash
MODULE   := $(shell sed -nr 's/^module ([a-z\-]+)$$/\1/p' go.mod)
GO_FILE  := $(MODULE).go
ifeq ($(output),)
OUT_FILE := dist/$(MODULE).bin
else
OUT_FILE := $(output)
endif

default:
	@{\
		set -e ;\
		os_release_id=$$(grep -E '^ID(_LIKE)?=' /etc/os-release | sort | head -n 1 | sed -r 's/ID(_LIKE)?=//' || true) ;\
		if [ "$$os_release_id" = "arch" ]; then \
			make --no-print-directory release type=dynamic ;\
		else \
			make --no-print-directory release type=static ;\
		fi ;\
	}

start:
	go run $(GO_FILE)

clean:
	go clean || true
	rm -rf dist/*
	rm -f "$(OUT_FILE)"

# Build development binary
.PHONY: build
build:
	go build -v -o $(OUT_FILE) $(GO_FILE)
	@ls -lh "$(OUT_FILE)"

release:
	@{\
		set -e ;\
		if [ "$(type)" != "static" ] && [ "$(type)" != "dynamic" ]; then \
			echo "The type parameter must be \"static\" or \"dynamic\"" ;\
			exit 1 ;\
		fi ;\
		echo "Building $(type) binary (GOARCH: $(GOARCH) GOARM: $(GOARM))..." ;\
		if [ -z "$${skip_clean}" ]; then make --no-print-directory clean; fi ;\
		export VERSION_FLAGS="-X main.Version=$$(make --no-print-directory version) -X main.CommitDate=$$(make --no-print-directory commit-date)" ;\
		case "$(type)" in \
			static) \
				make --no-print-directory release-static ;\
			;;\
			dynamic) \
				make --no-print-directory release-dynamic ;\
			;;\
		esac ;\
	}
	@make --no-print-directory postbuild

# Build statically linked production binary
release-static:
	@{\
		set -e ;\
		args=(-a -trimpath -mod=readonly -modcacherw -ldflags "-w -s $${VERSION_FLAGS}") ;\
		if [ "$${GOARCH}" != "arm" ]; then \
			args+=("-buildmode=pie") ;\
		fi ;\
		CGO_ENABLED=0 go build "$${args[@]}" -o "$(OUT_FILE)" $(GO_FILE) ;\
	}

# Build dinamically linked production binary
release-dynamic:
	@{\
		set -e ;\
		export CGO_CPPFLAGS="$${CPPFLAGS}" ;\
		export CGO_CFLAGS="$${CFLAGS}" ;\
		export CGO_CXXFLAGS="$${CXXFLAGS}" ;\
		export CGO_LDFLAGS="$${LDFLAGS}" ;\
		args=(-a -trimpath -mod=readonly -modcacherw -ldflags "-linkmode external -w -s $${VERSION_FLAGS}") ;\
		if [ "$${GOARCH}" != "arm" ]; then \
			args+=("-buildmode=pie") ;\
		fi ;\
		go build "$${args[@]}" -o "$(OUT_FILE)" $(GO_FILE) ;\
	}

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
				BIN_FILE="$(OUT_FILE)" \
				VERSION=$${VERSION} \
				RELEASE=$${RELEASE} \
				./scripts/mkdeb.py ;\
				rm -rf "$${PREFIX}" ;\
		fi ;\
	}
	@ls -lh "$(OUT_FILE)"*

docker:
	@{\
		set -e ;\
		image_name=$(MODULE)_$$(openssl rand -hex 8) ;\
		container_name=$(MODULE)_$$(openssl rand -hex 8) ;\
		docker build \
			--build-arg "CURRENT_TAG=$${CURRENT_TAG}" \
			--build-arg "COMMIT_DATE=$$(make --no-print-directory commit-date)" \
			--tag $${image_name} . ;\
		docker run -d --rm --name $${container_name} --entrypoint tail $${image_name} -f /dev/null ;\
		make --no-print-directory clean ;\
		docker cp $${container_name}:/go/src/app/dist . ;\
		docker stop -t 0 $${container_name} ;\
		docker rmi --no-prune $${image_name} ;\
	}

# Print the value of the VERSION variable if available, otherwise get version
# based on the latest git tag
version:
	@{\
		TAG=$$(git describe --tags 2>/dev/null | sed -r 's/^v([0-9]+\.[0-9]+\.[0-9]+).*/\1/') ;\
		VER=$${VERSION:-$${TAG}} ;\
		if [ ! -z "$$VER" ]; then \
			echo "$$VER" ;\
			exit 0 ;\
		fi ;\
		echo Unknown ;\
	}

# Print the value of the COMMIT_DATE variable if available, otherwise get commit
# date from the last git commit
commit-date:
	@{\
		set -e ;\
		if [ ! -z "$$COMMIT_DATE" ]; then \
			echo "$$COMMIT_DATE" ;\
			exit 0 ;\
		fi ;\
		git log -1 --no-merges --format=%cI ;\
	}

# Run unit tests on all packages
test:
	go test -v ./...

# Update Go version in go.mod file, keep in mind that -go must contain a major
# and a minor version number (i.e. not the last one)
update:
	go mod edit -go=$$(go version | sed -r 's/.*go([1-9]+\.[1-9]+)\..*/\1/')

install:
	install -Dm755 "$(OUT_FILE)" "$(PREFIX)/usr/bin/$(MODULE)"
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
