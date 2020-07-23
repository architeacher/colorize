test: bench integration profile race unit unit-short ## to setup the external used tools.

$(BENCH_TESTS_PATH):
	@if [ ! -d $(BENCH_TESTS_PATH) ] ; then mkdir -p $@ 2>&1 ; fi

$(BENCH_TESTS_PATH)/%.bin: $(BENCH_TESTS_PATH)
	@echo "$(INFO_CLR)$(MSG_PRFX) Benchmarking with cpu and memory profiles$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                                 \
		-bench=.                                \
		-benchmem                               \
		-count=$(BENCH_TESTS_COUNT)             \
		-cpuprofile $(BENCH_TESTS_PATH)/cpu.out \
		-memprofile $(BENCH_TESTS_PATH)/mem.out \
		-o $@				       		        \
		-parallel $(PARALLEL_TESTS)             \
		-run=NONE                               \
		-tags bench                             \
		-trace  $(BENCH_TESTS_PATH)/trace.out   \
		-timeout $(BENCH_TEST_TIMEOUT)          \
		$(GO_FLAGS) $(APP_DIR) 2>&1

bench: ## to run benchmark tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Benchmarking tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                        \
	    -bench .                       \
	    -benchmem 					   \
	    -count=$(BENCH_TESTS_COUNT)    \
	    -parallel $(PARALLEL_TESTS)    \
	    -run=NONE 					   \
	    -tags bench                    \
	    -timeout $(BENCH_TEST_TIMEOUT) \
	    $(GO_FLAGS) $(SRC_PKGS) 2>&1

bench-with-profile: $(BENCH_TESTS_PATH)/profile.bin ## to run benchmark with cpu and memory profiles.

clean-bench: ## to clean bench tests generated files.
	@echo "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up bench tests generated files$(MSG_SFX)$(NO_CLR)"
	@rm -rf "$(BENCH_TESTS_PATH)" 2>&1

clean-tests: clean-bench ## to clean coverage generated files.
	@echo "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up tests generated files$(MSG_SFX)$(NO_CLR)"

cpu-profile-serve: bench-with-profile ## to serve cpu benchmark profile over http - useful only if building remote/headless.
	@echo "$(INFO_CLR)$(MSG_PRFX) Serving cpu profile on port 8081$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool pprof -http=":8081" $(BENCH_TESTS_PATH)/profile.bin $(BENCH_TESTS_PATH)/cpu.out 2>&1

dump-assembly: ## to generate compiler assembly.
	@echo "$(INFO_CLR)$(MSG_PRFX) Generating compiler assembly$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool objdump $(BENCH_TESTS_PATH)/profile.bin > $(BENCH_TESTS_PATH)/profile.asm 2>&1

integration: build ## to run integration tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Integration tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
	    -cover                      \
	    -parallel $(PARALLEL_TESTS) \
	    -tags integration           \
	    -timeout $(TEST_TIMEOUT)    \
	    $(GO_FLAGS) ./... 2>&1

mem-profile-serve: bench-with-profile ## to serve memory benchmark profile over http - useful only if building remote/headless.
	@echo "$(INFO_CLR)$(MSG_PRFX) Serving memory profile on port 8082$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool pprof -http=":8082" $(BENCH_TESTS_PATH)/profile.bin $(BENCH_TESTS_PATH)/mem.out 2>&1

profile: profile-pdf profile-svg ## to get benchmark profiles in svg and pdf formats.

profile-%: bench-with-profile ## to get benchmark profiles in specified format.
	@$(eval FORMAT=$(firstword $(subst -, , $*)))
	@echo "$(INFO_CLR)$(MSG_PRFX) Generating bench profile in $(FORMAT) format$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool pprof -$(FORMAT) $(BENCH_TESTS_PATH)/profile.bin $(BENCH_TESTS_PATH)/cpu.out > $(BENCH_TESTS_PATH)/cpu.$(FORMAT) 2>&1
	@$(GO) tool pprof -$(FORMAT) $(BENCH_TESTS_PATH)/profile.bin $(BENCH_TESTS_PATH)/mem.out > $(BENCH_TESTS_PATH)/mem.$(FORMAT) 2>&1

race: ## to run long unit tests with race conditions detection coverage.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests with race cover$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                        \
	    -race                          \
	    -cpu=1,2,4                     \
	    -parallel $(PARALLEL_TESTS)    \
	    -tags $(GO_TAGS)               \
	    -timeout $(BENCH_TEST_TIMEOUT) \
	    $(GO_FLAGS) ./... 2>&1

trace-serve: bench-with-profile ## to serve runtime execution tracer.
	@echo "$(INFO_CLR)$(MSG_PRFX) Serving runtime execution tracer$(MSG_SFX)$(NO_CLR)"
	@$(GO) tool trace $(BENCH_TESTS_PATH)/trace.out 2>&1

unit: ## to run long unit tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
        -cover                      \
        -parallel $(PARALLEL_TESTS) \
        -tags $(GO_TAGS)            \
        -timeout=$(TEST_TIMEOUT)    \
	    $(GO_FLAGS) ./... 2>&1

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
unit-short: ## to run short unit tests.
	@echo "$(INFO_CLR)$(MSG_PRFX) Unit tests (short)$(MSG_SFX)$(NO_CLR)"
	@$(GO) test                     \
	    -cover                      \
	    -parallel $(PARALLEL_TESTS) \
	    -tags $(GO_TAGS)            \
	    -test.short                 \
	    -timeout=$(TEST_TIMEOUT)    \
	    $(GO_FLAGS) ./... 2>&1
