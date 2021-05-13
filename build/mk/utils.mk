define displayProjectLogo
    # http://patorjk.com/software/taag/#p=testall&f=Slant&t=PROJECT_NAME
    printf "$(1)"
    cat assets/logo.txt 2> /dev/null || echo $(PROJECT_NAME)
    printf "$(NO_CLR)\n"
endef

utils: clean format help list list-modules nuke

clean: clean-bin clean-coverage clean-version clean-tests ## to clean up all generated directories/files.

format: ## to format all Go files.
	printf "$(INFO_CLR)$(MSG_PRFX) Formatting code$(MSG_SFX)$(NO_CLR)\n"
	test -z "$$($(GO_FMT) -s -l -w $(GO_FLAGS) $(GO_FILES) 2>&1 | tee /dev/stderr)"

help: ## to get help about the targets.
	$(call displayProjectLogo,$(OK_CLR)) 2>&1
	printf "$(INFO_CLR)Please use \`make <target>\`, Available options for <target> are:$(NO_CLR)\n"
	awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z._-]+:.*?## .*$$/ {sub("\\\\n", sprintf("\n%22c"," "), $$2); printf "  $(STAR) $(HELP_CLR)%-28s$(NO_CLR) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort -u 2>&1
	printf "\n$(INFO_CLR)Useful variables:$(NO_CLR)\n"
	awk 'BEGIN { FS = "[:?]=" } /^## /{x = substr($$0, 4); getline; if (NF >= 2) printf "  $(PLUS) $(DISCLAIMER_CLR)%-29s$(NO_CLR) %s\n", $$1, x}' $(MAKEFILE_LIST) | sort -u 2>&1

list: ## to list all targets.
	awk -F':' '/^[a-z0-9][^$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /);for(i in A)printf "$(LIST_CLR)%-30s$(NO_CLR)\n", A[i]}' $(MAKEFILE_LIST) | sort -u 2>&1

list-modules: ## to list go modules.
	printf "$(WARN_CLR)$(MSG_PRFX) ℹ️ Installed Go modules$(MSG_SFX)$(NO_CLR)\n"
	$(GO) list -u -m all 2>&1

print-%: ## to print arbitrary variables use `print-VARNAME`
	echo "$(DISCLAIMER_CLR)$*$(NO_CLR) = $($*)"

nuke: clean ## to do clean up and enforce removing the corresponding installed archive or binary.
	printf "$(WARN_CLR)$(MSG_PRFX) 🧹 Cleaning up Go dependencies$(MSG_SFX)$(NO_CLR)\n"
	$(GO) clean -i -r -cache --modcache -testcache $(GO_FLAGS) ./... net 2>&1
