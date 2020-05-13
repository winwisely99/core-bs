
# This make file uses composition to keep things KISS and easy.
# In the boilerpalte make files dont do any includes, because you will create multi permutations of possibilities.

# includes
include ./boilerplate/help.mk
include ./boilerplate/go.mk


# examples of how to override in root make file
override GO_FSPATH = $(PWD)
override GO_BUILD_OUT_FSPATH = $(GOPATH)/bin
override GO_BUILD_OUT_ALL_FSPATH = $(PWD)/targets
override BS_ROOT_FSPATH = XXX
GO_ARCH=go-arch
override GO_ARCH=go-arch_override
override GO_BIN_NAME=bs

STATIK_DEST = $(PWD)/statiks

.PHONY: help this-statiks this-scan-statiks this-build

## Print all settings
this-print:
	$(MAKE) go-print

## Example to Print Variable override from make
this-print-ex:
	# prints specific overides
	@echo BS_ROOT_FSPATH: 	$(BS_ROOT_FSPATH)

	@echo GO_ARCH: 	$(GO_ARCH)
	
## Example to Print Variable override from env
this-print-env-ex:
	# prints specific overides
	# Example call of override from env variable
	# ``` GO_ARCH=GO_ARCH_FROMENV make -e print-env ```
	@echo BS_ROOT_FSPATH: 	$(BS_ROOT_FSPATH)

	@echo GO_ARCH: 	$(GO_ARCH)


## Build this.
this-build:
	$(MAKE) go-build

## Build for all arch and platforms
this-build-all:
	$(MAKE) go-build-all

## Delete all of the builds
this-build-clean-all:
	$(MAKE) go-build-clean-all

## Delete the build.
this-build-clean:
	$(MAKE) go-build-clean
	# delete all generated stuff
	rm -rf $(PWD)/statiks

#this-statiks:
#	go run sdk/cmd/bp/main.go -t $(PWD)/.tmp -o $(PWD)/statiks

