# Disable all the default make stuff
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:

ARTIFACT_DIR ?= artifacts

# Add more binaries to the BINARIES variable. These should exist as cmd/*/main.go
BINARIES ?=

# The platform specific binaries to build. Add more here if you need to support windows or other architectures.
ARTIFACT_BINARIES := $(addprefix $(ARTIFACT_DIR)/,$(BINARIES:=.linux.amd64) $(BINARIES:=.darwin.amd64) $(BINARIES:=.windows.amd64))

GOFLAGS := -mod=vendor
GO_BINARY := $(if $(USE_DOCKER),docker run --rm -t -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) -e GOFLAGS=$(GOFLAGS) -v $(PWD):/__ -w /__ golang:1.11 go,GOOS=$(GOOS) GOARCH=$(GOARCH) GOFLAGS=$(GOFLAGS) go)
GO_FILES := $(shell find $(subst /__/,$(PWD)/,$(shell $(GO_BINARY) list -f '{{ .Dir }}' ./...) -type f -name *.go))
GO_DEPENDENCIES := $(GO_FILES) go.sum go.mod

# Logging function used for some consistency - the if/else is used to turn off colors when no TTY exists
ifeq "$(shell [ -t 0 ] && echo 1)" "1"
LOG = @echo "\033[38;5;13m$(shell date):\033[0m"
else
LOG = @echo $(shell date):
endif

# Shasums command
SHASUMS := $(if $(uname | grep Darwin),shasum -a 256,sha256sum)

## Display a list of the documented make targets
.PHONY: help
help:
	@echo Documented Make targets:
	@perl -e 'undef $$/; while (<>) { while ($$_ =~ /## (.*?)(?:\n# .*)*\n.PHONY:\s+(\S+).*/mg) { printf "\033[36m%-30s\033[0m %s\n", $$2, $$1 } }' $(MAKEFILE_LIST) | sort

# ------------------------------------------------------------------------------
# NON-PHONY TARGETS
# ------------------------------------------------------------------------------

.PHONY: .FORCE
.FORCE:

# Create the artifactory directory
$(ARTIFACT_DIR):
	$(LOG) Creating artifact dir $@..
	@mkdir -p $@

# Build the development binary
$(BINARIES): $(GO_DEPENDENCIES)
	$(LOG) Building $@...
	@$(GO_BINARY) build -o $@ -v ./cmd/$(notdir $@)

# Write the current git describe state
$(ARTIFACT_DIR)/describe: .FORCE | $(ARTIFACT_DIR)
	@git describe --always --tags --long --dirty --match 'v[0-9]*.[0-9]*.[0-9]*' 2>/dev/null > $@.new
	@cmp $@.new $@ 2>/dev/null || mv -f $@.new $@
	@rm -f $@.new

# Write the semantic version string file
$(ARTIFACT_DIR)/semver: describe = $(shell cat $(ARTIFACT_DIR)/describe)
$(ARTIFACT_DIR)/semver: $(ARTIFACT_DIR)/describe
	$(LOG) Writing new semver file $@..
	@echo $(if $(describe:v%=),0.0.1-0-)$(describe) | sed -e 's/^v//; s/-/\#/; s/-/./g; s/\#/-/' > $@

# Build the release binaries
$(ARTIFACT_BINARIES): t = $(subst ., ,$(notdir $@))
$(ARTIFACT_BINARIES): GOOS = $(word 2,$(t))
$(ARTIFACT_BINARIES): GOARCH = $(word 3,$(t))
$(ARTIFACT_BINARIES): $(GO_DEPENDENCIES) $(ARTIFACT_DIR)/semver
	$(LOG) Building $@ with version $(shell cat $(ARTIFACT_DIR)/semver)..
	@$(GO_BINARY) build \
		-o $@ \
		-ldflags '-X main.version=$(shell cat $(ARTIFACT_DIR)/semver)+$(USER) -X "main.buildDate=$(shell date)"' \
		-v \
		./cmd/$(word 1,$(t))

# Write a shasums file for the artifact binaries
$(ARTIFACT_DIR)/SHASUMS: $(ARTIFACT_BINARIES)
	$(LOG) Calculating shasums for binaries..
	@cd $(dir $@) && $(SHASUMS) $(notdir $(ARTIFACT_BINARIES)) > $(notdir $@)

# Write coverage file
$(ARTIFACT_DIR)/coverage.out: $(GO_DEPENDENCIES) | $(ARTIFACT_DIR)
	$(LOG) Running unittests with coverage..
	@$(GO_BINARY) test -v -cover -covermode=count -coverprofile=$@ ./...
	$(LOG) Running vet checks..
	@$(GO_BINARY) vet ./...

# Build a docker image
$(ARTIFACT_DIR)/docker-%: tag = $*:$(shell cat $(ARTIFACT_DIR)/semver).$(shell date '+%Y%m%d%H%M%S')
$(ARTIFACT_DIR)/docker-%: $(ARTIFACT_DIR)/%.linux.amd64 $(ARTIFACT_DIR)/semver .dockerignore Dockerfile
	$(LOG) Building Docker image..
	@docker build -t $(tag) .
	@echo $(tag) > $@

# ------------------------------------------------------------------------------
# PHONY TARGETS
# ------------------------------------------------------------------------------

## Build development binary on local system
.PHONY: dev
dev: $(BINARIES)
	$(LOG) Development binaries available at: $^

.PHONY: version
version: $(ARTIFACT_DIR)/semver
	@cat $<

## Build release binaries
.PHONY: build
build: $(ARTIFACT_BINARIES) $(ARTIFACT_DIR)/SHASUMS
	$(LOG) Artifact binaries available at: $^

## Run unittests
.PHONY: test
test: $(ARTIFACT_DIR)/coverage.out

## Print function coverage
.PHONY: coverage
coverage: $(ARTIFACT_DIR)/coverage.out
	$(LOG) Printing function Coverage..
	@$(GO_BINARY) tool cover -func $<

## Open coverage in browser
.PHONY: coverage-html
coverage-html: $(ARTIFACT_DIR)/coverage.out
	$(LOG) Printing function Coverage..
	@$(GO_BINARY) tool cover -html $<

## Build docker image
.PHONY: image
image: $(ARTIFACT_DIR)/docker-$(firstword $(BINARIES))
	$(LOG) Image is available as $(shell cat $<)

## Clean old files and artifacts
.PHONY: clean
clean:
	$(LOG) Cleaning old files..
	@rm -rfv $(ARTIFACT_DIR) $(BINARIES)
