############## VARIABLES SECTION ##############
# The binary to build.
BIN := kubeswap
VERSION := $(shell cat internal/version/version.go | grep Version | cut -d= -f2 | tr -d '[[:space:]]' | sed s/\"//g)

# The input main.go
MAIN = main.go

# The output binary
OUTBIN = dist/$(OS)_$(ARCH)/$(BIN)

############## TARGETS SECTION ##############
.PHONY: all test clean
build: # @HELP builds for current GOOS/GOARCH
build:
	@goreleaser build --snapshot --single-target --skip-validate

snapshot: # @HELP generate a snapshot for all OS_ARCH combinations
snapshot:
	@goreleaser --snapshot --skip-publish

release: # @HELP releases a new version for all OS_ARCH combinations
release:
	@goreleaser release

version: # @HELP outputs the version string
version:
	@echo $(VERSION)

dep-upgrade: # @HELP upgrades all dependencies
dep-upgrade:
	@go get -u ./...
	@go mod tidy

clean: # @HELP removes built binaries and temporary files
clean:
	@rm -rf dist

test: # @HELP executes the test/test.sh script
test:
	@./test/test.sh

help: # @HELP prints this message
help:
	@echo "VARIABLES:"
	@echo "  BIN = $(BIN)"
	@echo "  OS = $(OS)"
	@echo "  ARCH = $(ARCH)"
	@echo "  VERSION = $(VERSION)"
	@echo
	@echo "TARGETS:"
	@grep -E '^.*: *# *@HELP' Makefile            \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-30s %s\n", $$1, $$2 };  \
	    '
