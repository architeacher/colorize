.DEFAULT: default
.DEFAULT_GOAL := default
# Delete targets on nonzero exit status.
.DELETE_ON_ERROR: ;
.EXPORT_ALL_VARIABLES: ;
.NOTPARALLEL: ;
.ONESHELL: ;
.SHELLFLAGS = -c

# Unless ECHO_RECIPES is defined we do not echo the line of the recipes before they are executed.
ifneq (true, ${ECHO_RECIPES})
.SILENT: ;
endif

# An empty, phony .FORCE target, that will cause external targets to always be built, so that the
# Makefile there will handle dependencies.
# https://gist.github.com/mschubert/a0e4f3aeaf3558431890
.PHONY: .FORCE
.FORCE: ;

GO := go
GO_FLAGS ?= $(GO_FLAGS:)
GO_FILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*" 2> /dev/null)

# Directories which hold app source (not vendored), e.g. cmd, pkg, ...etc.
SRC_PKGS ?= $(shell $(GO) list ./... 2> /dev/null | grep -v "/vendor/" |  grep -v "/cmd")

# include .env file without overriding vars setup in the console.
ifneq (,$(wildcard ${DOT_ENV_FILE}))
  VARS = $(shell sed -ne 's/ *\#.*$$//; /./ s/=/?=/ p' ${DOT_ENV_FILE})
  $(foreach v,${VARS},$(eval export ${v}))
endif

include build/mk/config/settings.mk
include build/mk/*.mk

##@ Please use `make <target>`, Available options for <target> are:
.PHONY: all
all: all-build coverage profile test utils validate ## to run all targets.

.PHONY: default
default: help ## to run default target.
