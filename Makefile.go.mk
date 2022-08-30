# In short to add a binary you usually want to add it either to
# the BINARIES variable if it is a linux Binary and WINDOWS_BINARIES
# if it is a Windows Binary (There is also DARWIN_BINARIES but only for migctl)
#
# The build flags (including injecting the build version) are part of the build/gobuild.sh
# script and the build/report_build_info.sh script
#
#
# List of all binaries to build
BINARIES:=\
	./cmd/containerdbg \
	./cmd/test-binary \
	./cmd/scale-binary \
	$(NULL)

WINDOWS_BINARIES:=\
	./cmd/containerdbg \
	$(NULL)

DARWIN_BINARIES:=\
	./cmd/containerdbg \
	$(NULL)

CGO_ENABLED_BINARIES:=\
	$(NULL)

# Global Variables
SHELL := /bin/bash -o pipefail

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

CLANG ?= clang-13
CFLAGS := -O2 -g -Wall $(CFLAGS) -Werror

# locations where artifacts are stored
# cumulatively track the directories/files to delete after a clean
DIRS_TO_CLEAN:=
FILES_TO_CLEAN:=

# If GOPATH is not set by the env, set it to a sane value
GOPATH ?= $(shell go env GOPATH)
export GOPATH

# Note that disabling cgo here adversely affects go get.  Instead we'll rely on this
# to be handled in build/gobuild.sh
# export CGO_ENABLED=0

# It's more concise to use GO?=$(shell which go)
# but the following approach uses a more efficient "simply expanded" :=
# variable instead of a "recursively expanded" =
ifeq ($(origin GO), undefined)
  GO:=$(shell which go)
endif
ifeq ($(GO),)
  $(error Could not find 'go' in path.  Please install go, or if already installed either add it to your path or set GO to point to its directory)
endif

LOCAL_ARCH := $(shell uname -m)
ifeq ($(LOCAL_ARCH),x86_64)
GOARCH_LOCAL := amd64
else
GOARCH_LOCAL := $(LOCAL_ARCH)
endif
export GOARCH ?= $(GOARCH_LOCAL)

LOCAL_OS := $(shell uname)
ifeq ($(LOCAL_OS),Linux)
   export GOOS_LOCAL = linux
else ifeq ($(LOCAL_OS),Darwin)
   export GOOS_LOCAL = darwin
else
   $(error "This system's OS $(LOCAL_OS) isn't recognized/supported")
   # export GOOS_LOCAL ?= windows
endif

export GOOS ?= $(GOOS_LOCAL)

# Invoke make VERBOSE=1 to enable echoing of the command being executed
export VERBOSE ?= 0

BUILDTYPE_DIR := release

# @todo allow user to run for a single $PKG only?
PACKAGES_CMD := GOPATH=$(GOPATH) $(GO) list ./...
GO_FILES_CMD := find . -name '*.go' | grep -v -E '$(GO_EXCLUDE)'

# Environment for tests, the directory containing istio and deps binaries.
# Typically same as GOPATH/bin, so tests work seemlessly with IDEs.

# Using same package structure as pkg/
MY_MAKEFILE := $(abspath $(lastword $(MAKEFILE_LIST)))
MY_MK_DIR:= $(dir $(MY_MAKEFILE))
export OUT_DIR:=$(MY_MK_DIR)out
export OUT_LINUX := $(OUT_DIR)/linux_amd64/$(BUILDTYPE_DIR)
export OUT_WINDOWS := $(OUT_DIR)/windows_amd64/$(BUILDTYPE_DIR)
export OUT_DARWIN := $(OUT_DIR)/darwin_amd64/$(BUILDTYPE_DIR)
export REPO_ROOT := $(shell git rev-parse --show-toplevel)

.PHONY: default
default: depend build test

# The point of these is to allow scripts to query where artifacts
# are stored so that tests and other consumers of the build don't
# need to be updated to follow the changes in the Makefiles.
# Note that the query needs to pass the same types of parameters
# (e.g., DEBUG=0, GOOS=linux) as the actual build for the query
# to provide an accurate result.
.PHONY: where-is-out
where-is-out:
	@echo ${OUT_DIR}

.PHONY: depend init

init:
	@mkdir -p ${OUT_DIR}/logs

OUTPUT_DIRS := \
	$(OUT_LINUX) \
	$(OUT_WINDOWS) \
	$(OUT_DARWIN)

generate: export BPF_CLANG := $(CLANG)
generate: export BPF_CFLAGS := $(CFLAGS)
generate:
	go generate ./...

depend: $(all_proto_go) generate init buildinfo deployment | $(OUTPUT_DIRS)

DIRS_TO_CLEAN := $(OUTPUT_DIRS)

$(OUTPUT_DIRS):
	@mkdir -p $@

.PHONY: precommit format lint buildcache

precommit: format lint

#TBD
lint:
	@echo TBD Lint

# Build with -i to store the build caches into $GOPATH/pkg
buildcache:
	GOBUILDFLAGS=-i $(MAKE) -f Makefile build-go-linux


.PHONY: build-go
build-go: depend build-go-linux build-go-windows build-go-darwin

.PHONY: buildinfo
buildinfo:
	@./build/report_build_info.sh > ./build/.buildinfo

.PHONY: go-build-deps
go-build-deps: $(all_proto_go)

SHARED_BUILD_FLAGS := STATIC=0 GOARCH=amd64 LDFLAGS='-extldflags -static -s -w' 
# TODO: fix ugly dependancy on common_interfaces
.PHONY: build-go-linux
build-go-linux: depend go-build-deps
	GOOS=linux $(SHARED_BUILD_FLAGS) $(REPO_ROOT)/build/gobuild.sh $(OUT_LINUX)/ $(BINARIES)

.PHONY: build-go-windows
build-go-windows: depend $(all_proto_go) $(OUT_WINDOWS)
	GOOS=windows $(SHARED_BUILD_FLAGS) $(REPO_ROOT)/build/gobuild.sh $(OUT_LINUX)/ $(WINDOWS_BINARIES)

.PHONY: build-go-darwin
build-go-darwin: depend $(all_proto_go) $(OUT_DARWIN)
	GOOS=darwin $(SHARED_BUILD_FLAGS) $(REPO_ROOT)/build/gobuild.sh $(OUT_LINUX)/ $(DARWIN_BINARIES)

.PHONY: clean
go_clean:
	$(RM) $(FILES_TO_CLEAN)
	$(RM) -r $(DIRS_TO_CLEAN)

CLEAN_TARGETS+=go_clean
