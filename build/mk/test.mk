.PHONY: test
test: integration race unit unit-short ## to setup the external used tools.

.PHONY: integration
integration: build ## to run integration tests.
	$(call printMessage,"🤝 Integration tests",$(INFO_CLR))
	$(GO) test \
		-parallel $(PARALLEL_TESTS) \
		-tags integration \
		-timeout $(TEST_TIMEOUT) \
		$(GO_FLAGS) \
		$(APP_DIR)/... 2>&1

.PHONY: race
race: ## to run long unit tests with race conditions detection.
	$(call printMessage,"🏁 Running unit tests (+ race detector)",$(INFO_CLR))
	CGO_ENABLED=1 \
	$(GO) test \
		-cpu=1,2,4 \
		-parallel $(PARALLEL_TESTS) \
		-race \
		-tags $(GO_TAGS) \
		-timeout $(BENCH_PROFILER_TIMEOUT) \
		$(GO_FLAGS) \
		$(APP_DIR)/... 2>&1

.PHONY: unit
unit: ## to run long unit tests.
	$(call printMessage,"🏃 Running unit tests",$(INFO_CLR))
	$(GO) test \
		-parallel $(PARALLEL_TESTS) \
		-tags $(GO_TAGS) \
		-timeout=$(TEST_TIMEOUT) \
		$(GO_FLAGS) \
		$(APP_DIR)/... 2>&1

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
.PHONY: unit-short
unit-short: ## to run short unit tests.
	$(call printMessage,"🏃 Running unit tests (short)",$(INFO_CLR))
	$(GO) test \
		-parallel $(PARALLEL_TESTS) \
		-tags $(GO_TAGS) \
		-test.short \
		-timeout=$(TEST_TIMEOUT) \
		$(GO_FLAGS) \
		$(APP_DIR)/... 2>&1
