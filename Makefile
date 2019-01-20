BINARIES := md-slides
ARTIFACT_DIR := artifacts
include base.Makefile

$(ARTIFACT_DIR)/pages:
	@mkdir -pv $@

$(ARTIFACT_DIR)/pages/index.html: md-slides README.md | $(ARTIFACT_DIR)/pages
	$(LOG) Generating $@..
	@./md-slides html README.md $(dir $@)
	@cp -v windmill.jpeg $(dir $@)

## Generate pages pdf content
.PHONY: pages
pages: $(ARTIFACT_DIR)/pages/index.html
	$(LOG) Pages available at $(dir $<)

## Run more detailed integration tests
.PHONY: integ-test
integ-test: md-slides | $(ARTIFACT_DIR)

	$(LOG) testing pdf subcommand
	@./md-slides pdf -tmp-dir artifacts README.md artifacts/test.pdf

	$(LOG) testing install.sh script
	INSTALL_DIRECTORY=$(shell mktemp -d -p $(ARTIFACT_DIR)) OVERRIDE_RELEASE_TAG=v0.7.0 \
		./install.sh
