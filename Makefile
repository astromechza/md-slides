BINARIES := md-slides
ARTIFACT_DIR := artifacts
include base.Makefile

$(ARTIFACT_DIR)/pages:
	@mkdir -pv $@

$(ARTIFACT_DIR)/pages/index.html: md-slides SLIDES.md | $(ARTIFACT_DIR)/pages
	$(LOG) Generating $@..
	@./md-slides html --source SLIDES.md --target-dir $(dir $@)

## Generate pages pdf content
.PHONY: pages
pages: $(ARTIFACT_DIR)/pages/index.html
	$(LOG) Pages available at $(dir $<)
