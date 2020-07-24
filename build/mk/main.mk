GO := go
GO_FLAGS ?= $(GO_FLAGS:)
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*" 2> /dev/null)

# Directories which hold app source (not vendored), e.g. cmd, pkg, ...etc.
SRC_PKGS ?= $(shell $(GO) list ./... 2> /dev/null | grep -v "/vendor/" |  grep -v "/cmd")

include build/mk/utils.mk
include build/mk/build.mk
include build/mk/coverage.mk
include build/mk/test.mk
include build/mk/validate.mk

all: all-build all-coverage utils test validate ## to run all targets.
default: build
