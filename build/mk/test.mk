test: bench integration profile race unit unit-short ## to setup the external used tools.

$(BENCH_TESTS_PATH):
	@if [ ! -d $(BENCH_TESTS_PATH) ] ; then mkdir -p $(BENCH_TESTS_PATH) 2>&1 ; fi

bench: ## to run benchmark tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Benchmarking tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                               \
	    -count=$(BENCH_TESTS_COUNT)           \
	    -run=NONE -bench .                    \
	    -benchmem -parallel $(PARALLEL_TESTS) \
	    -timeout $(TEST_TIMEOUT)              \
	    -tags bench                           \
	    $(GO_FLAGS) $(SRC_PKGS) 2>&1

clean-bench: ## to clean bench tests generated files.
	@echo "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up bench tests generated files$(MSG_SFX)$(NO_CLR)"
	@rm -rf "$(BENCH_TESTS_PATH)" 2>&1

clean-tests: clean-bench ## to clean coverage generated files.
	@echo "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up tests generated files$(MSG_SFX)$(NO_CLR)"

cpu-profile-serve: profile ## to serve cpu bench mark profiles over http - useful only if building remote/headless.
	@echo "$(INFO_CLR)$(MSG_PRFX) Serving memory profile on port 8081$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool pprof -http=":8081" $(BENCH_TESTS_PATH)/cpu.out $(BENCH_TESTS_PATH)/cpu.svg 2>&1

integration: build ## to run integration tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Integration tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
	    -cover                      \
	    -parallel $(PARALLEL_TESTS) \
	    -timeout $(TEST_TIMEOUT)    \
	    -tags integration           \
	    $(GO_FLAGS) ./... 2>&1

mem-profile-serve: profile ## to serve memory bench mark profiles over http - useful only if building remote/headless.
	@echo "$(INFO_CLR)$(MSG_PRFX) Serving memory profile on port 8082$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool pprof -http=":8082" $(BENCH_TESTS_PATH)/mem.out $(BENCH_TESTS_PATH)/mem.svg 2>&1
	@rm -rf "$(TESTS_PATH)" 2>&1

profile: $(BENCH_TESTS_PATH) ## to get bench mark profiles.
	@echo "$(INFO_CLR)$(MSG_PRFX) Bench tests check$(MSG_SFX)$(NO_CLR)"
	$(GO) test                                  \
	    -count=$(BENCH_TESTS_COUNT)             \
	    -run=NONE                               \
	    -bench=.                                \
	    -benchmem                               \
	    -o $(BENCH_TESTS_PATH)/test.bin         \
	    -cpuprofile $(BENCH_TESTS_PATH)/cpu.out \
	    -memprofile $(BENCH_TESTS_PATH)/mem.out \
	    -parallel $(PARALLEL_TESTS)             \
	    -timeout $(TEST_TIMEOUT)                \
	    -tags bench                             \
	    $(GO_FLAGS) ./$(APP_DIR) 2>&1
	@$(GO) tool pprof --svg $(BENCH_TESTS_PATH)/test.bin $(BENCH_TESTS_PATH)/cpu.out > $(BENCH_TESTS_PATH)/cpu.svg 2>&1
	@$(GO) tool pprof --svg $(BENCH_TESTS_PATH)/test.bin $(BENCH_TESTS_PATH)/mem.out > $(BENCH_TESTS_PATH)/mem.svg 2>&1

race: ## to run long unit tests with race conditions detection coverage.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests with race cover$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
	    -race                       \
	    -cpu=1,2,4                  \
	    -parallel $(PARALLEL_TESTS) \
	    -timeout $(TEST_TIMEOUT)    \
	    -tags $(GO_TAGS)            \
	    $(GO_FLAGS) ./... 2>&1

unit: ## to run long unit tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
        -cover                      \
        -parallel $(PARALLEL_TESTS) \
        -timeout=$(TEST_TIMEOUT)    \
        -tags $(GO_TAGS)            \
	    $(GO_FLAGS) ./... 2>&1

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
unit-short: ## to run short unit tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests (short)$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
	    -test.short                 \
	    -cover                      \
	    -parallel $(PARALLEL_TESTS) \
	    -timeout=$(TEST_TIMEOUT)    \
	    -tags $(GO_TAGS)            \
	    $(GO_FLAGS) ./... 2>&1
