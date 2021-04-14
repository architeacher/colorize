test: bench integration profile race unit unit-short ## to setup the external used tools.

$(BENCH_TESTS_PATH):
	if [ ! -d $(BENCH_TESTS_PATH) ] ; then mkdir -p $@ 2>&1 ; fi

$(BENCH_PROFILE): $(BENCH_TESTS_PATH)
	printf "$(INFO_CLR)$(MSG_PRFX) Benchmarking with cpu and memory profiles$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                              \
		-bench=.                            \
		-benchmem                           \
		-count=$(BENCH_TESTS_COUNT)         \
		-cpuprofile $(BENCH_CPU_PROFILE) 	\
		-memprofile $(BENCH_MEMORY_PROFILE) \
		-o $@				       		    \
		-parallel $(PARALLEL_TESTS)         \
		-run=NONE                           \
		-tags bench                         \
		-timeout $(BENCH_TEST_TIMEOUT)      \
		$(GO_FLAGS) $(APP_DIR)/. 2>&1

$(TRACE_TESTS_PATH):
	if [ ! -d $(TRACE_TESTS_PATH) ] ; then mkdir -p $@ 2>&1 ; fi

bench: ## to run benchmark tests.
	printf "$(INFO_CLR)$(MSG_PRFX) Benchmarking tests$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                         \
	    -bench .                       \
	    -benchmem 					   \
	    -count=$(BENCH_TESTS_COUNT)    \
	    -parallel $(PARALLEL_TESTS)    \
	    -run=NONE 					   \
	    -tags bench                    \
	    -timeout $(BENCH_TEST_TIMEOUT) \
	    $(GO_FLAGS) $(APP_DIR)/... 2>&1

bench-with-profile: $(BENCH_PROFILE) ## to run benchmark with cpu and memory profiles.

clean-bench: ## to clean bench tests generated files.
	printf "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up bench tests generated files$(MSG_SFX)$(NO_CLR)\n"
	rm -rf "$(BENCH_TESTS_PATH)" 2>&1

clean-tests: clean-bench clean-trace ## to clean coverage generated files.
	printf "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up tests generated files$(MSG_SFX)$(NO_CLR)\n"
	rm -rf "$(TESTS_PATH)" 2>&1

clean-trace: ## to clean trace profile generated files.
	printf "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up trace profile generated files$(MSG_SFX)$(NO_CLR)\n"
	rm -rf "$(TRACE_TESTS_PATH)" 2>&1

cpu-profile-serve: bench-with-profile ## to serve cpu benchmark profile over http - useful only if building remote/headless.
	printf "$(INFO_CLR)$(MSG_PRFX) Serving cpu profile on port 8081$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool pprof -http=":8081" $(BENCH_PROFILE) $(BENCH_CPU_PROFILE) 2>&1

dump-assembly: ## to generate compiler assembly.
	printf "$(INFO_CLR)$(MSG_PRFX) Generating compiler assembly$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool objdump $(BENCH_PROFILE) > $(BENCH_TESTS_PATH)/profile.asm 2>&1

integration: build ## to run integration tests.
	printf "$(INFO_CLR)$(MSG_PRFX) Integration tests$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                      \
	    -parallel $(PARALLEL_TESTS) \
	    -tags integration           \
	    -timeout $(TEST_TIMEOUT)    \
	    $(GO_FLAGS) $(APP_DIR)/... 2>&1

mem-profile-serve: bench-with-profile ## to serve memory benchmark profile over http - useful only if building remote/headless.
	printf "$(INFO_CLR)$(MSG_PRFX) Serving memory profile on port 8082$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool pprof -http=":8082" $(BENCH_PROFILE) $(BENCH_MEMORY_PROFILE) 2>&1

profile: profile-pdf profile-png profile-svg ## to get benchmark profiles in pdf, png and svg formats.

profile-%: bench-with-profile ## to get benchmark profiles in specified format.
	$(eval FORMAT=$(firstword $(subst -, , $*)))
	printf "$(INFO_CLR)$(MSG_PRFX) Generating bench profile in $(FORMAT) format$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool pprof -$(FORMAT) $(BENCH_PROFILE) $(BENCH_CPU_PROFILE) > $(BENCH_TESTS_PATH)/cpu.$(FORMAT) 2>&1
	$(GO) tool pprof -$(FORMAT) $(BENCH_PROFILE) $(BENCH_MEMORY_PROFILE) > $(BENCH_TESTS_PATH)/mem.$(FORMAT) 2>&1

race: ## to run long unit tests with race conditions detection coverage.
	printf "$(INFO_CLR)$(MSG_PRFX) Unit tests with race cover$(MSG_SFX)$(NO_CLR)\n"
	CGO_ENABLED=1                      \
	$(GO) test                         \
	    -race                          \
	    -cpu=1,2,4                     \
	    -parallel $(PARALLEL_TESTS)    \
	    -tags $(GO_TAGS)               \
	    -timeout $(BENCH_TEST_TIMEOUT) \
	    $(GO_FLAGS) $(APP_DIR)/... 2>&1

trace: $(TRACE_TESTS_PATH) ## to run tests with tracing options.
	printf "$(INFO_CLR)$(MSG_PRFX) Tests with trace options$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                               \
		-benchmem                            \
		-blockprofile $(BLOCK_TRACE_PROFILE) \
		-mutexprofile $(MUTEX_TRACE_PROFILE) \
		-parallel $(PARALLEL_TESTS)          \
		-run=NONE                            \
		-tags trace                          \
		-trace $(EXEC_TRACE_PROFILE)    	 \
		-timeout $(TEST_TIMEOUT)       		 \
		$(GO_FLAGS) $(APP_DIR)/. 2>&1

trace-serve: trace ## to serve runtime execution tracer.
	printf "$(INFO_CLR)$(MSG_PRFX) Serving runtime execution tracer$(MSG_SFX)$(NO_CLR)\n"
	$(GO) tool trace $(MUTEX_TRACE_PROFILE) 2>&1

unit: ## to run long unit tests.
	printf "$(INFO_CLR)$(MSG_PRFX) Unit tests$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                      \
        -parallel $(PARALLEL_TESTS) \
        -tags $(GO_TAGS)            \
        -timeout=$(TEST_TIMEOUT)    \
	    $(GO_FLAGS) $(APP_DIR)/... 2>&1

# Quick test. You can bypass long tests using: `if testing.Short() { t.Skip("Skipping in short mode.") }`.
unit-short: ## to run short unit tests.
	printf "$(INFO_CLR)$(MSG_PRFX) Unit tests (short)$(MSG_SFX)$(NO_CLR)\n"
	$(GO) test                      \
	    -parallel $(PARALLEL_TESTS) \
	    -tags $(GO_TAGS)            \
	    -test.short                 \
	    -timeout=$(TEST_TIMEOUT)    \
	    $(GO_FLAGS) $(APP_DIR)/... 2>&1
