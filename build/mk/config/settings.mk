# See https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
NO_CLR = \033[0m
AZURE = \x1b[1;38;5;45m
CYAN = \x1b[96m
GREEN = \x1b[1;38;5;113m
OLIVE = \x1b[1;36;5;113m
MAGENTA = \x1b[38;5;170m
ORANGE =  \x1b[1;38;5;208m
RED = \x1b[91m
YELLOW = \x1b[1;38;5;227m

INFO_CLR := ${AZURE}
DISCLAIMER_CLR := ${MAGENTA}
ERROR_CLR := ${RED}
HELP_CLR := ${CYAN}
OK_CLR := ${GREEN}
ITEM_CLR := ${OLIVE}
LIST_CLR := ${ORANGE}
WARN_CLR := ${YELLOW}

STAR := ${ITEM_CLR}[${NO_CLR}${LIST_CLR}*${NO_CLR}${ITEM_CLR}]${NO_CLR}
PLUS := ${ITEM_CLR}[${NO_CLR}${WARN_CLR}+${NO_CLR}${ITEM_CLR}]${NO_CLR}

MSG_PRFX := ==>
MSG_SFX := ...

DEPENDENCIES := github.com/client9/misspell/cmd/misspell \
                github.com/mattn/goveralls \
                github.com/wadey/gocovmerge

## Path to .env file.
DOT_ENV_FILE ?= $(CURDIR)/.env

## To echo recipes, you can do "make ECHO_RECIPES=true".
ECHO_RECIPES ?= false

## To disable root, you can do "make SUDO=".
SUDO ?= $(shell echo "sudo -E" 2> /dev/null)

## Should do cross compile for other OSs/Architectures or not.
CROSS_BUILD ?= false

# https://github.com/golang/go/blob/master/src/go/build/syslist.go
## List of possible platforms for cross compile.
TARGET_PLATFORMS ?=linux darwin

# amd64 (x86-64), i386 (x86 or x86-32), arm64 (AArch64), arm (ARM), ppc64le (IBM POWER8), s390x (IBM System z) ...etc.
## List of possible architectures for cross compile.
TARGET_ARCHS ?=amd64 i386 arm64 arm ppc64le s390x

## Operating system to build for.
OS ?= $(shell uname -s 2>&1 | tr '[:upper:]' '[:lower:]')

## Architecture to build for.
ARCH ?= amd64

## Extra flags to pass to 'go' when building.
GO_FLAGS ?=

## Version file path.
VERSION_FILE ?= $(CURDIR)/.version

## Current version.
VERSION ?= $(shell cat $(VERSION_FILE) 2> /dev/null || git describe --match 'v[0-9]*' --abbrev=0 2> /dev/null || echo NA)

## If true, disable optimizations and does NOT strip the binary.
DEBUG ?= false

## If true, "build" will produce a static binary (cross compile always produce static build regardless).
STATIC ?= true

## Base application directory name.
APP_DIR_NAME ?= .

## Path where the application files are located at.
APP_DIR ?= $(CURDIR)/$(APP_DIR_NAME)

## Path where the main Go file is located at.
CMD_DIR ?= $(APP_DIR)/cmd/

## Path where the Go packages are located at.
PKG_DIR ?= $(APP_DIR)/pkg

## Path where Go package version directory is located at.
PKG_VERSION_DIR ?= $(PKG_DIR)/version

## Go package version template file name.
PKG_VERSION_TEMPLATE ?= version-template.go.dist

## Set an output prefix, which is the local directory if not specified.
ARTIFACTS_DIR ?= ${CURDIR}/.artifacts

## Set the binary directory, where the compiled should go to.
BINARY_DIR ?= ${ARTIFACTS_DIR}/bin

## Set the binary file name prefix.
BINARY_PREFIX ?= $(shell basename $(CURDIR) 2> /dev/null)

## Base path used to install.
INSTALLATION_BASE_PATH ?= /usr/local/bin

## Go generated files.
GO_GENERATED_DIR = ${ARTIFACTS_DIR}/go

## 🐳 Should use containers locally to do build, lint, ... etc.
CONTAINERIZE ?= true

# Tests
## Set tests path.
TESTS_PATH ?= $(GO_GENERATED_DIR)/tests

## The number of parallel tests.
PARALLEL_TESTS ?= 8

## Test timeout.
TEST_TIMEOUT ?= 18s

## Test time multiplier flag name.
TEST_TIME_MULTIPLIER_FLAG ?= timeMultiplier

## Test time multiplier value.
TEST_TIME_MULTIPLIER ?= 1

## Set linter results path.
LINTER_PROFILE_PATH ?= $(TESTS_PATH)/linter
LINTER_PROFILE := $(LINTER_PROFILE_PATH)/gl-code-quality-report.json

# Coverage tests
## Set coverage mode {atomic|count|set}.
COVERAGE_MODE ?= atomic

## Set coverage profile path.
COVERAGE_PATH ?= $(TESTS_PATH)/coverage
COVERAGE_PROFILE := $(COVERAGE_PATH)/profile.out
COVERAGE_HTML := $(COVERAGE_PATH)/index.html
COVERAGE_TEXT := $(COVERAGE_PATH)/function.out

# Profiling
## Set profiler path.
PROFILER_PATH ?= $(GO_GENERATED_DIR)/profiler

## Bench memory profile path.
BENCH_CPU_PROFILE ?= $(PROFILER_PATH)/mem.out
## Bench cpu profile path.
BENCH_MEMORY_PROFILE ?= $(PROFILER_PATH)/cpu.out
## Bench binary profiler path.
BENCH_PROFILE ?= $(PROFILER_PATH)/bench.bin

## The times that each test profile and benchmark would run.
BENCH_PROFILE_COUNT ?= 3

## Bench tests timeout.
BENCH_PROFILER_TIMEOUT ?= 18m

# Goroutine blocking profile path.
BLOCK_TRACE_PROFILE ?= $(PROFILER_PATH)/blocking.out
# Mutex contention profile path.
MUTEX_TRACE_PROFILE ?= $(PROFILER_PATH)/mutex.out
# Execution trace profile path.
EXEC_TRACE_PROFILE ?= $(PROFILER_PATH)/trace.out
