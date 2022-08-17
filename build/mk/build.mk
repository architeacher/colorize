GO_OS ?= $(shell $(GO) env GOOS 2> /dev/null)
GO_ARCH ?= $(shell $(GO) env GOARCH 2> /dev/null)
GO_VERSION := $(shell $(GO) version | grep -o -E "go[0-9\.]+" 2> /dev/null)

# Disable usage of CGO.
CGO_ENABLED := 0

ifeq ($(ARCH),)
    ARCH := $(GO_ARCH)
endif

ifeq ($(OS),)
    OS := $(GO_OS)
endif

# Valid OS and Architecture target combinations. Check https://golang.org/doc/install/source#environment
# Or run `go tool dist list -json | jq`
VALID_OS_ARCH := "[darwin/arm64][darwin/amd64][linux/amd64][linux/386][linux/arm64][linux/arm][linux/ppc64le][openbsd/amd64][windows/amd64]"

os.darwin := Darwin
os.linux := Linux
os.openbsd := OpenBSD
os.windows := Windows

arch.amd64 := x86_64
arch.386 := i386
arch.arm64 := aarch64
arch.arm := armhf
arch.ppc64le := ppc64le
arch.s390x := s390x

BINARY_NAME := "${BINARY_PREFIX}-${os.$(OS)}-${arch.$(ARCH)}"
TARGET_BINARY := "${BINARY_DIR}/${BINARY_NAME}"

# Package main path.
PKG_NAMESPACE ?= $(shell $(GO) list -e ./ 2> /dev/null)

# Get the current local branch name from git (if we can, this may be blank).
GIT_BRANCH := $(shell git symbolic-ref --short HEAD 2> /dev/null || git rev-parse --abbrev-ref HEAD 2> /dev/null)
GIT_COMMIT := $(shell git rev-list -1 HEAD 2> /dev/null)
# Get the git commit.
GIT_DIRTY := $(shell test -n "`git status --porcelain --untracked-files=no 2> /dev/null`" && echo "+CHANGES" || true 2> /dev/null)

# Build Flags
# The default version that's chosen when pushing the images. Can/should be overridden.
BUILD_VERSION ?= $(shell git describe --always --abbrev=8 --tags --dirty='-Changes' 2> /dev/null | cut -d "v" -f 2 2> /dev/null)
BUILD_HASH ?= git-$(shell git rev-parse --short=18 HEAD 2> /dev/null)
BUILD_TIME ?= $(shell date +%FT%T%z 2> /dev/null)

# If we don't set the build version it defaults to dev.
ifeq ($(BUILD_VERSION),)
	BUILD_VERSION := $(shell cat $(CURDIR)/.version 2> /dev/null || echo dev)
endif

BUILD_ENV ?= $(BUILD_ENV:)

GO_BUILD_FLAGS ?= -a -installsuffix cgo

EXTLD_FLAGS ?=

# Check if we are not building for darwin, and honoring static.
IS_DARWIN_HOST ?= $(shell echo $(GO_OS) | egrep -i -c "darwin" 2> /dev/null)
IS_STATIC ?= $(shell echo $(STATIC) | egrep -i -c "true" 2>&1)

# Below, we are building a boolean circuit that says "$(IS_DARWIN_HOST) && $(IS_STATIC)"
ifeq ($(shell echo $$(( $(IS_DARWIN_HOST) * $(IS_STATIC) )) 2> /dev/null), 0)
# The flags we are passing to go build. -extldflags -static for making a static binary,
# or -linkmode external for linking external C libraries into the binary.
    override EXTLD_FLAGS +=-lm -static -lstdc++ -lpthread -static-libstdc++
endif

# -X version.BuildHash for telling the Go binary which build hash is used in this version,
# -X version.BuildTime for telling the Go binary the build time,
# -X version.GitBranch for telling the Go binary the git branch used,
# -X version.GitCommit for telling the Go binary the git commit used,
# -X version.GoVersion for telling the Go binary the go version used,
# -X main.version for telling the Go binary which version it is.
GO_LINKER_FLAGS ?=-s \
        -v \
        -w \
        -X ${PKG_NAMESPACE}/version.BuildHash=$(BUILD_HASH) \
        -X ${PKG_NAMESPACE}/core.BuildTime=$(BUILD_TIME) \
        -X ${PKG_NAMESPACE}/core.GitBranch=$(GIT_BRANCH) \
        -X ${PKG_NAMESPACE}/core.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) \
        -X ${PKG_NAMESPACE}/core.GoVersion=$(GO_VERSION) \
        -X main.version=$(BUILD_VERSION)

ifdef EXTLD_FLAGS
    GO_LINKER_FLAGS	+= -extldflags "$(EXTLD_FLAGS)"
endif

GO_GC_FLAGS :=-trimpath=$(CURDIR)

# Honor debug
ifeq ($(DEBUG), true)
	# Disable function inlining and variable registration.
	GO_GC_FLAGS +=all=-N -l
endif

GO_ASM_FLAGS :=-trimpath=$(CURDIR)

# netgo for enforcing the native Go DNS resolver.
GO_TAGS ?= netgo

GO_ENV_FLAGS ?= $(GO_ENV_FLAGS:)
GO_ENV_FLAGS += $(BUILD_ENV)

extension = $(patsubst windows, .exe, $(filter windows, $(1)))

define goCross
	$(if $(findstring [$(1)/$(2)],$(VALID_OS_ARCH)), \
		$(call printMessage,"🏗️ Building binary for [$(1)/$(2)]",$(OK_CLR)); \
		$(eval CMD := CGO_ENABLED=$(CGO_ENABLED) $\
									GOOS=$(1) $\
									GOARCH=$(2) $\
									$(GO_ENV_FLAGS) $\
									$(GO) build $\
									$(GO_BUILD_FLAGS) $\
									-asmflags='$(GO_ASM_FLAGS)' $\
									-gcflags='$(GO_GC_FLAGS)' $\
									-ldflags '$(GO_LINKER_FLAGS)' $\
									-o $(BINARY_DIR)/$(BINARY_PREFIX)-${os.$(1)}-${arch.$(2)}$(call extension,$(GO_OS)) $\
									-tags $(GO_TAGS) $\
									$(GO_FLAGS) $\
									$(CMD_DIR)) \
		printf "${INFO_CLR}	$${CMD} ${NO_CLR}\n"; \
		$$(eval $${CMD}) 2>&1, \
		printf "${ERROR_CLR}No defined build target for: \"[${1}/${2}]\"${NO_CLR}\n $\
					${INFO_CLR}Defined build targets are: ${VALID_OS_ARCH}.${NO_CLR}\n" \
	)
endef

define buildTargets
	$(foreach GO_OS, $(3), $(foreach GO_ARCH, $(4), $(call $(1), $(2)$(GO_OS), -$(GO_ARCH))))
endef

define getDependency
	GO111MODULE=off $(GO) get -u -v $(GO_FLAGS) $(1) 2>&1;
endef

define replaceInFile
    $(if $(findstring $(IS_DARWIN_HOST),1),  \
        sed -i '' "s|$(1)|$(2)|g" $(3) 2>&1, \
        sed -i -e "s|$(1)|$(2)|g" $(3) 2>&1)
endef

.PHONY: all-build
all-build: build-bin build-version build-x clean-bin clean-version deps go-generate go-install install update-pkg-version uninstall version

.PHONY: build
ifneq ($(CONTAINERIZE), true)

build: build-bin ## to install dependencies and build out a binary.

else

build: deploy

endif

.PHONY: build-bin
build-bin: ## to build out a binary.
	$(if $(filter $(CROSS_BUILD), true), \
	    $(MAKE) build-x, \
	    $(MAKE) $(call buildTargets, addprefix, build-bin-for-, $(OS), $(ARCH)))

.PHONY: build-bin-for-%
build-bin-for-%:
	$(eval TARGET_PLATFORM=$(firstword $(subst -, , $*)))
	$(eval TARGET_ARCH=$(or $(word 2,$(subst -, , $*)),$(value 2)))
	$(call goCross,$(TARGET_PLATFORM),$(TARGET_ARCH))

.PHONY: build-version
build-version:  ## to get the current build version.
	echo $(BUILD_VERSION)

.PHONY: build-x
build-x: $(shell find . -type f -name '*.go') ## to build for cross platforms.
#	$(foreach GO_OS, $(TARGET_PLATFORMS), $(foreach GO_ARCH, $(TARGET_ARCHS), $(call goCross,$(GO_OS),$(GO_ARCH))))
	$(foreach GO_OS, $(TARGET_PLATFORMS), $(foreach GO_ARCH, $(TARGET_ARCHS), $(MAKE) $(call buildTargets, addprefix, build-bin-for-, $(GO_OS), $(GO_ARCH))))

.PHONY: clean-bin
clean-bin: ## to clean up generated binaries.
	$(call printMessage,"🧹 Cleaning up generated binaries",$(WARN_CLR))
	rm -rf "${BINARY_DIR}" 2>&1

.PHONY: clean-version
clean-version: ## to remove generated version file.
	$(call printMessage,"🧹 Cleaning up generated version file",$(WARN_CLR))
	rm -f  "${VERSION_FILE}"

.PHONY: deps
deps: ## to get required dependencies.
	$(call printMessage,"⬇️ Installing required dependencies",$(OK_CLR))
	$(foreach dependency, $(DEPENDENCIES), $(call getDependency,$(dependency)))

.PHONY: go-generate
go-generate: ## to generate Go related files.
	$(call printMessage,"Generating files via Go generate",$(OK_CLR))
	$(GO) generate $(GO_FLAGS) $(SRC_PKGS) 2>&1

.PHONY: go-install
go-install: update-pkg-version ## to install the Go related/dependent commands and packages.
	$(call printMessage,"⬇️ Installing Go related dependencies",$(INFO_CLR))
	$(GO) install \
	    -ldflags "-X $(PKG_NAMESPACE)/$(APP_DIR_NAME)/pkg/version.VERSION=${VERSION}" \
	    -installsuffix "static" \
	    -tags $(GO_TAGS) \
	    -v $(GO_FLAGS) \
	    $(SRC_PKGS) 2>&1

.PHONY: install
install: ## to install the generated binary.
	$(call printMessage,"⬇️ Installing generated binary",$(INFO_CLR))
	if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	sudo cp $(TARGET_BINARY) $(INSTALLATION_BASE_PATH) 2>&1

.PHONY: kill
kill: ## to send a kill signal to the running process of the binary.
	$(call printMessage,"🥷️ Sending kill signal",$(WARN_CLR))
	pkill "${args}" $(notdir $(TARGET_BINARY)) > /dev/null 2>&1

.PHONY: mods-download
mods-download: ## to download all required modules.
	$(call printMessage,"⬇️ Installing required modules",$(INFO_CLR))
	$(GO) mod download $(GO_FLAGS) 2>&1

.PHONY: mods-vendor
mods-vendor: ## to download all required modules as vendor.
	$(call printMessage,"⬇️ Installing required modules as vendor",$(INFO_CLR))
	$(GO) mod vendor $(GO_FLAGS) 2>&1

.PHONY: run
run: ## to run the generated binary, and build a new one if not existed.
	$(call printMessage,"🏃 Running generated binary",$(OK_CLR))
	if [ ! -f $(TARGET_BINARY) ] ; then $(MAKE) build; fi
	$(TARGET_BINARY) "${args}" 2>&1

.PHONY: uninstall
uninstall: ## to uninstall generated binary.
	$(call printMessage,"Uninstalling generated binary",$(WARN_CLR))
	sudo rm -rf $(INSTALLATION_BASE_PATH)/$(BINARY_NAME) 2>&1

.PHONY: update-pkg-version
update-pkg-version: ## to update package version.
	$(call printMessage,"🔃 Updating Go package version",$(INFO_CLR))
  ifneq ($(wildcard $(PKG_VERSION_DIR)/$(PKG_VERSION_TEMPLATE)),)
		cp $(PKG_VERSION_DIR)/$(PKG_VERSION_TEMPLATE) $(PKG_VERSION_DIR)/version.go 2>&1
		echo $(VERSION) > $(VERSION_FILE) 2>&1
		$(call replaceInFile,{{VERSION}},$(VERSION),$(PKG_VERSION_DIR)/version.go)
  endif

.PHONY: version
version:  ## to get the current version.
	$(call printMessage,"Current tagged version",$(INFO_CLR))
	$(eval version_commit := $(shell git rev-parse --short=8 $(VERSION) 2> /dev/null || echo NA))
	$(call printMessage,"${VERSION} -> ${version_commit}",$(OK_CLR))
