define displayProjectLogo
    # http://patorjk.com/software/taag/#p=testall&f=Slant&t=PROJECT_NAME
    printf "${1}"
    cat assets/logo.txt 2> /dev/null || echo $(PROJECT_NAME)
    printf "$(NO_CLR)\n"
endef

define printMessage
    printf "${2}$(MSG_PRFX) %s$(MSG_SFX)$(NO_CLR)\n" ${1} 2>&1
endef

$(ARTIFACTS_DIR):
	if [ ! -d $(ARTIFACTS_DIR) ] ; then mkdir -p $(ARTIFACTS_DIR) 2>&1 ; fi

.PHONY: utils
utils: clean help list list-modules nuke size stats todo

.PHONY: clean
clean: clean-artifacts clean-version ## to clean up all generated directories/files.

.PHONY: clean-artifacts
clean-artifacts: clean-bin clean-coverage clean-linting clean-profiling ## to clean up generated artifacts.
	$(call printMessage,"🧹 Cleaning up generated artifacts",$(WARN_CLR))
	rm -rf "${ARTIFACTS_DIR}" 2>&1

.PHONY: help
help: ## to get help about the targets.
	$(call displayProjectLogo,$(OK_CLR)) 2>&1
	awk 'BEGIN {FS = ":.*?## "}; \
		/^[a-zA-Z._-]+%?:.*?## .*$$/ {sub("\\\\n", sprintf("\n%22c"," "), $$2); \
		printf "  $(STAR) $(HELP_CLR)%-28s$(NO_CLR) %s\n", $$1, $$2} \
		/^##@/ { printf "\n$(INFO_CLR)%s$(NO_CLR)\n", substr($$0, 5) } \
    /^##-/ { printf "  %-17s\n", substr($$0, 5) }' \
		$(MAKEFILE_LIST) | sort -u 2>&1
	printf "\n$(INFO_CLR)Useful variables:$(NO_CLR)\n"
	awk 'BEGIN { FS = "[:?]=" }; \
		/^## /{x = substr($$0, 4); getline; \
		if (NF >= 2) printf "  $(PLUS) $(DISCLAIMER_CLR)%-29s$(NO_CLR) %s\n", $$1, x}' \
		$(MAKEFILE_LIST) | sort -u 2>&1

.PHONY: list
list: ## to list all targets.
	awk -F':' '/^[a-z0-9][^$#\/\t=]*:([^=]|$$)/ {split($$1,A,/ /); \
		for(i in A)printf "$(LIST_CLR)%-30s$(NO_CLR)\n", A[i]}' \
		$(MAKEFILE_LIST) | sort -u 2>&1

.PHONY: list-modules
list-modules: ## to list go modules.
	$(call printMessage,"ℹ️ Installed Go modules",$(INFO_CLR))
	$(GO) list -m -u -mod=mod all 2>&1

.PHONY: print-%
print-%: ## to print arbitrary variables use `print-VARNAME`
	echo "$(DISCLAIMER_CLR)$*$(NO_CLR) = $($*)"

.PHONY: nuke
nuke: clean ## to do clean up and enforce removing the corresponding installed archive or binary.
	$(call printMessage,"☢️ Cleaning up Go dependencies",$(WARN_CLR))
	$(GO) clean \
		-i \
		-r \
		-cache \
		--modcache \
		-testcache \
		$(GO_FLAGS) \
		$(APP_DIR)/... net 2>&1

.PHONY: size
size: ## to show size of imports.
	$(call printMessage,"📐 Calculating size of the imports",$(INFO_CLR))
	eval `go build -work -a 2>&1` && \
		find "$${WORK}" -type f -name "*.a" | \
		xargs -I{} du -hxs "{}" | \
		sort -rh | \
		sed -e "s:${WORK}/::g" 2>&1

.PHONY: stats
stats: ## to output source statistics.
	$(call printMessage,"📊 Calculating source statistics",$(INFO_CLR))
	cloc --exclude-dir=node_modules,vendor . 2>&1

.PHONY: todo
todo: $(ARTIFACTS_DIR) ## to output to-do items per file.
	$(call printMessage,"🔎️ Searching for todos",$(INFO_CLR))
	todos="$$(todoPrefix="TODO"; \
		grep \
		--color \
		--exclude-dir=./.artifacts \
		--exclude-dir=./assets \
		--exclude-dir=./vendor \
		--exclude-dir=./node_modules \
		--text \
		-inRo \
		" $${todoPrefix}:.*" . )" ; \
	if [ -n "$${todos}" ]; then \
		echo "${ITEM_CLR}$${todos}${NO_CLR}"; \
		echo "$${todos}" > "${ARTIFACTS_DIR}/TODOs.txt"; \
	fi
