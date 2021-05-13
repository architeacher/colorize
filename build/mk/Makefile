include build/mk/config.mk
include build/mk/main.mk

.DEFAULT: default
.DEFAULT_GOAL := default
# Delete targets on nonzero exit status
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
