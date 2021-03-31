COVERAGE_PACKAGES := $(SRC_PKGS)

define coverPackage
	$(eval FILE_NAME=$(shell basename $(1) 2>&1))
	$(GO) test                                   		 \
	    -cover                                   		 \
	    -covermode=$(COVERAGE_MODE)              		 \
	    -coverprofile $(COVERAGE_PATH)/$(FILE_NAME).part \
	    -parallel $(PARALLEL_TESTS)              		 \
	    -tags $(GO_TAGS)            	 		 		 \
	    -timeout $(TEST_TIMEOUT)                 		 \
	    $(GO_FLAGS)                              		 \
	    $(1) 2>&1;
endef

# Goveralls binary.
GOVERALLS_BIN := $(GOPATH)/bin/goveralls
GOVERALLS := $(shell [ -x $(GOVERALLS_BIN) ] && echo $(GOVERALLS_BIN) || echo '' 2> /dev/null)
coverage: cover coverage-browse

$(COVERAGE_HTML): $(COVERAGE_PROFILE)
	printf "$(INFO_CLR)$(MSG_PRFX) Coverage HTML export$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool cover -html="$(COVERAGE_PROFILE)" -o "$(COVERAGE_HTML)" $(GO_FLAGS) 2>&1

$(COVERAGE_PROFILE):
	if [ ! -d $(COVERAGE_PATH) ] ; then mkdir -p $(COVERAGE_PATH) 2>&1 ; fi
	$(foreach package, $(COVERAGE_PACKAGES), $(call coverPackage,$(package)))

	echo "mode: $(COVERAGE_MODE)" > $(COVERAGE_PROFILE) 2>&1
	# tail -q -n +2 $(COVERAGE_PATH)/*.part
	find $(COVERAGE_PATH) -type f -name '*.part' -exec grep -h -v "mode: $(COVERAGE_MODE)" {} + >> "$(COVERAGE_PROFILE)" 2>&1

clean-coverage: ## to clean coverage generated files.
	printf "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up coverage generated files$(MSG_SFX)$(NO_CLR)\n"
	rm -rf "$(COVERAGE_PATH)" 2>&1

cover: update-pkg-version $(COVERAGE_PROFILE) ## to run test with coverage and report that out to profile.
	printf "$(INFO_CLR)$(MSG_PRFX) Coverage check$(MSG_SFX)$(NO_CLR)\n"

coverage-browse: coverage-html ## to browse coverage results in html format.
	printf "$(INFO_CLR)$(MSG_PRFX) Opening browser$(MSG_SFX)$(NO_CLR)\n"
	open "$(COVERAGE_HTML)" 2>&1

coverage-html: $(COVERAGE_HTML) ## to export coverage results to html format "$(COVERAGE_PATH)/index.html".
	printf "$(INFO_CLR)$(MSG_PRFX) Generating HTML coverage to $(COVERAGE_PATH)/index.html$(MSG_SFX)$(NO_CLR)\n"

coverage-send: $(COVERAGE_PROFILE) ## to send the results to coveralls.
	printf "$(INFO_CLR)$(MSG_PRFX) Sending coverage$(MSG_SFX)$(NO_CLR)\n"
	$(if $(GOVERALLS), , $(error Please run make get-deps))
	$(GOVERALLS) -service travis-ci -coverprofile="$(COVERAGE_PROFILE)" -repotoken $(COVERALLS_TOKEN) 2>&1

coverage-serve: coverage-html ## to serve coverage results over http - useful only if building remote/headless.
	printf "$(INFO_CLR)$(MSG_PRFX) Serving coverage on port 8888$(MSG_SFX)$(NO_CLR)\n"
	cd "$(COVERAGE_PATH)" && python -m SimpleHTTPServer 8888 2>&1
