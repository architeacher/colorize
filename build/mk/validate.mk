GO_LINT ?= golangci-lint
GO_FMT ?= gofmt

.PHONY: validate
validate: format-check lint misspell vet ## to validate the code.

$(LINTER_PROFILE_PATH):
	if [ ! -d $(LINTER_PROFILE_PATH) ] ; then mkdir -p $(LINTER_PROFILE_PATH) 2>&1 ; fi

.PHONY: clean-linting
clean-linting: ## to clean generated linting files.
	$(call printMessage,"🧹 Cleaning up generated linting files",$(WARN_CLR))
	rm -rf "$(LINTER_PROFILE_PATH)" 2>&1

.PHONY: format-check
format-check: ## to check if the Go files are formatted correctly.
	$(call printMessage,"🏃 Running code format check",$(INFO_CLR))
	diff="$$($(GO_FMT) -d -e -l -s ${args} $(GO_FILES) 2>&1)"; \
	if [ -n "$${diff}" ]; then \
		echo "$(ERROR_CLR)Error, Please run 'make format-check args=-w' and commit the result:$(NO_CLR)"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: lint
lint: $(LINTER_PROFILE_PATH) ## to run linter against Go files.
	$(call printMessage,"🏃 Running linter",$(INFO_CLR))
	$(eval CMD ?= "$(GO_LINT) run ${args} --out-format code-climate | jq 'map(.severity = (.severity // \"\"))' | tee $(LINTER_PROFILE) | jq -r '.[] | \"\(.location.path):\(.location.lines.begin) \(.description)\"'")
	$(if $(filter $(CONTAINERIZE), true), \
			${MAKE} lint-dockerized CMD=${CMD}, \
	    $$(eval ${CMD}) 2>&1)

ifeq ($(CONTAINERIZE), true)

.PHONY: lint-dockerized
lint-dockerized: ## 🐳 to run linter against Go files in docker container.
	$(call printMessage,"🏃 Running linter in 🐳 container",$(INFO_CLR))
	${DOCKER} run \
		--name golangci-lint \
		--rm \
		-v $(CURDIR):/app \
		-w /app \
		golangci/golangci-lint:latest-alpine \
		$(CMD) 2>&1

endif

.PHONY: misspell
misspell: ## to run misspell.
	$(call printMessage,"🏃 Running misspell",$(INFO_CLR))
	diff="$$(find . \
					-type d \( -path ./.git -o -path ./${ARTIFACTS_DIR} -o -path ./.idea -o -path ./vendor \) \
					-o \
					-prune \
					-type f \
					-exec misspell "${args}" -error "${0}" {} /dev/null \;)"; \
	if [ -n "$${diff}" ]; then \
		echo "${ERROR_CLR}Error, Please run 'make misspell args=-w' and commit the result:${NO_CLR}"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
.PHONY: vet
vet: ## to run detection on dead code.
	$(call printMessage,"🏃 Running vet",$(INFO_CLR))
	test -z "$$($(GO) vet $(GO_FLAGS) $(SRC_PKGS) 2>&1 | tee /dev/stderr)" 2>&1
