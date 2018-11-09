VERSION := $(shell git describe --always --tags --dirty)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
GIT_DATE := $(shell git log -1 --format=%cI)

ALL_FILES := $(shell find . -type f -name '*.go')
PKG := github.com/AstromechZA/md-slides/cmd/mdslides
BINARY := md-slides

BUILD_OS := $(shell uname | tr "[:upper:]" "[:lower:]")
ifeq "$(BUILD_OS)" "linux"
SHACMD := sha256sum
else
SHACMD := shasum -a 256
endif

# make the main binary
.PHONY: dev
dev: $(BINARY)

# tasks for creating dist directory
dist/:
	mkdir dist

# tasks for creating binaries
DISTRIBUTABLES = $(BINARY) dist/$(BINARY).linux.amd64 dist/$(BINARY).darwin.amd64 dist/$(BINARY).windows.amd64
GOOS = $(shell echo $@ | grep amd64 | rev | cut -d. -f 2 | rev)
GOARCH = $(shell echo $@ | grep amd64 | rev | cut -d. -f 1 | rev)

# official distributable version
$(DISTRIBUTABLES): dist/ $(ALL_FILES)
	@echo Building $@..
	@CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-o $@ \
		-ldflags "-X main.gitHash=$(GIT_COMMIT) -X main.gitDate=$(GIT_DATE) -X main.version=$(VERSION)" \
		$(PKG)

# shasums file for dist
dist/SHA256SUMS: $(DISTRIBUTABLES)
	cd dist && $(SHACMD) md-slides.* > SHA256SUMS

pages/:
	mkdir pages

pages/index.html: pages/ $(BINARY) README.md
	@./$(BINARY) html README.md pages/
	@cp -v windmill.jpeg pages/

.PHONY: test
test: $(BINARY)
	@./$(BINARY) -version

# release build
.PHONY: release
release: dist/SHA256SUMS $(DISTRIBUTABLES)
	ls -l $(DISTRIBUTABLES)
	cat dist/SHA256SUMS

# clean will just remove the main bits and pieces
.PHONY: clean
clean:
	@rm -rfv $(DISTRIBUTABLES)
	@rm -rfv dist
	@rm -rfv pages/index.html
