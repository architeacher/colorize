COVERAGE_PACKAGES := $(SRC_PKGS)

# Goveralls binary.
GOVERALLS_BIN := $(GOPATH)/bin/goveralls
GOVERALLS := $(shell [ -x $(GOVERALLS_BIN) ] && echo $(GOVERALLS_BIN) || echo '' 2> /dev/null)

define coverPackage
	$(eval FILE_NAME=$(shell basename $(1) 2>&1))
	$(GO) test \
		-cover \
		-covermode=$(COVERAGE_MODE) \
		-coverprofile $(COVERAGE_PATH)/$(FILE_NAME).part \
		-parallel $(PARALLEL_TESTS) \
		-tags $(GO_TAGS) \
		-timeout $(TEST_TIMEOUT) \
		$(GO_FLAGS) \
		$(1) 2>&1;
endef

.PHONY: coverage
coverage: cover coverage-browse

$(COVERAGE_PROFILE):
	if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) 2>&1 ; fi
	$(foreach package, $(COVERAGE_PACKAGES), $(call coverPackage,$(package)))

	echo "mode: ${COVERAGE_MODE}" > $(COVERAGE_PROFILE) 2>&1
	# tail -q -n +2 $(COVERAGE_PATH)/*.part
	find $(COVERAGE_PATH) -type f -name '*.part' -exec grep -h -v "mode: ${COVERAGE_MODE}" {} + >> "${COVERAGE_PROFILE}" 2>&1

$(COVERAGE_HTML): $(COVERAGE_PROFILE)
	$(call printMessage,"Coverage HTML export",$(INFO_CLR))
	$(GO) tool cover -html="${COVERAGE_PROFILE}" -o "${COVERAGE_HTML}" ${args} 2>&1

$(COVERAGE_TEXT): $(COVERAGE_PROFILE)
	$(call printMessage,"Coverage Text export",$(INFO_CLR))
	$(GO) tool cover -func="${COVERAGE_PROFILE}" -o "${COVERAGE_TEXT}" ${args} 2>&1

.PHONY: clean-coverage
clean-coverage: ## to clean generated coverage files.
	$(call printMessage,"🧹 Cleaning up generated coverage files",$(WARN_CLR))
	rm -rf "${COVERAGE_PATH}" 2>&1

.PHONY: cover
cover: $(COVERAGE_PROFILE) ## to run test with coverage and report that out to profile.
	$(call printMessage,"Coverage check",$(INFO_CLR))

.PHONY: coverage-browse
coverage-browse: coverage-html ## to browse coverage results in html format.
	$(call printMessage,"Opening browser",$(INFO_CLR))
	open "${COVERAGE_HTML}" 2>&1

.PHONY: coverage-html
coverage-html: $(COVERAGE_HTML) ## to export coverage percentages to html format "${COVERAGE_HTML}".
	$(call printMessage,"Generating HTML coverage to ${COVERAGE_HTML}",$(INFO_CLR))

.PHONY: coverage-send
coverage-send: $(COVERAGE_PROFILE) ## to send the results to coveralls.
	$(call printMessage,"Sending coverage",$(INFO_CLR))
	$(if $(GOVERALLS), , $(error Please run make deps))
	$(GOVERALLS) -service travis-ci -coverprofile="${COVERAGE_PROFILE}" -repotoken $(COVERALLS_TOKEN) 2>&1

.PHONY: coverage-serve
coverage-serve: coverage-html ## to serve coverage percentages over http - useful only if building remote/headless.
	$(call printMessage,"Serving coverage on port 8888",$(INFO_CLR))
	cd "${COVERAGE_PATH}" && python -m SimpleHTTPServer 8888 2>&1

.PHONY: coverage-text
coverage-text: $(COVERAGE_TEXT) ## to export coverage percentages for each function to text format "${COVERAGE_TEXT}".
	$(call printMessage,"Generating TEXT coverage to ${COVERAGE_TEXT}",$(INFO_CLR))
