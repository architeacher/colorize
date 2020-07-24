GO_LINT ?= golint
GO_FMT ?= gofmt

validate: ineffassign format-check lint misspell vet ## to validate the code.

ineffassign: ## to run ineffassign.
	printf "$(INFO_CLR)$(MSG_PRFX) 🏃 Running ineffassign$(MSG_SFX)$(NO_CLR)\n"
	test -z "$$(ineffassign . | grep -v vendor/ | grep -v ".pb.go:" | tee /dev/stderr)"

format-check: ## to check if the Go files are formatted correctly.
	printf "$(INFO_CLR)$(MSG_PRFX) Checking code format$(MSG_SFX)$(NO_CLR)\n"
	diff=$$($(GO_FMT) -d -s $(GO_FILES)); 						\
	if [ -n "$$diff" ]; then                                    \
		echo "Please run 'make format' and commit the result:"; \
		echo "$${diff}";                                        \
		exit 1;                                                 \
	fi;

lint: ## to run linter against Go files.
	printf "$(INFO_CLR)$(MSG_PRFX) 🏃 Running linter$(MSG_SFX)$(NO_CLR)\n"
	$(GO_LINT) $(GO_FLAGS) $(SRC_PKGS)

misspell: ## to run misspell.
	printf "$(INFO_CLR)$(MSG_PRFX) 🏃 Running misspell$(MSG_SFX)$(NO_CLR)\n"
	test -z "$$(find . -type f | grep -v vendor/ | grep -v bin/ | grep -v .git/ | grep -v MAINTAINERS | xargs misspell | tee /dev/stderr)"

# Simplified dead code detector. Used for skipping certain checks on unreachable code
# (for instance, shift checks on arch-specific code).
# https://golang.org/cmd/vet/
vet: ## to run detection on dead code.
	printf "$(INFO_CLR)$(MSG_PRFX) 🏃 Running vet$(MSG_SFX)$(NO_CLR)\n"
	test -z "$$($(GO) vet $(GO_FLAGS) $(SRC_PKGS) 2>&1 | tee /dev/stderr)"
